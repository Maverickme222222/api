// Package main ...
package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Maverickme222222/api/handlers"
	"github.com/Maverickme222222/api/services"
)

const (
	// LogStrKeyModule is for use with the logger as a key to specify the module name.
	LogStrKeyModule = "module"
	// LogStrKeyService is for use with the logger as a key to specify the service name.
	LogStrKeyService = "service"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "TESTING"

func main() {

	z := zerolog.New(os.Stderr).With().Str(LogStrKeyService, "api").Timestamp().Logger()
	mainLog := z.With().Str(LogStrKeyModule, "main").Logger()
	mainLog.Info().Msg("starting server...")

	if err := run(&mainLog); err != nil {
		mainLog.Info().Msgf("main: error %s:", err.Error())
		os.Exit(1)
	}
}

func run(log *zerolog.Logger) error {
	log.Info().Msg("Welcome to the API Service :)")

	// =========================================================================
	// Configuration

	// Call env.Load func to ensure a .env file is loaded when available.
	// This command only affects dev environments.
	LoadDevEnv()

	var cfg struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"env:RELAY_APIHOST,default:localhost"`
			Port            int           `conf:"env:RELAY_PORT,default:9090"`
			ReadTimeout     time.Duration `conf:"env:RELAY_API_READ_TIMEOUT,default:5s"`
			WriteTimeout    time.Duration `conf:"env:RELAY_API_WRITE_TIMEOUT,default:5s"`
			ShutdownTimeout time.Duration `conf:"env:RELAY_API_SHUTDOWN_TIMEOUT,default:5s"`
		}
		Services struct {
			Users  string `conf:"env:USERS_SVC"`
			Emails string `conf:"env:EMAILS_SVC"`
		}
	}
	cfg.Version.Build = build
	cfg.Version.Desc = "API Service"

	// =========================================================================
	// App Starting

	expvar.NewString("build").Set(build)
	log.Info().Msgf("Started: Application initializing: version %q", build)
	defer log.Info().Msg("Completed")

	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			log.Info().Msg(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Info().Msgf("Config:\n%v\n", out)

	// =========================================================================
	// Register GRPC Client Services

	log.Info().Msg("main: Initializing Support for Downstream Services")

	ctx := context.Background()

	registeredServices, err := services.Register(ctx,
		cfg.Services.Users,
		cfg.Services.Emails)

	if err != nil {
		return errors.Wrap(err, "could not connect to one or more services")
	}

	// Start API Service

	log.Info().Msg("main: Initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		ReadHeaderTimeout: time.Second * 30,
		Addr:              fmt.Sprintf("%s:%d", cfg.Web.APIHost, cfg.Web.Port),
		Handler:           handlers.API(build, registeredServices, log),
		ReadTimeout:       cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Info().Msgf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Info().Msgf("main: %v: Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			if err := api.Close(); err != nil {
				return errors.Wrap(err, "could not stop server gracefully")
			}
			return errors.Wrap(err, "could not stop server")
		}
	}

	return nil
}

// LoadDevEnv loads .env file if present
func LoadDevEnv() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
