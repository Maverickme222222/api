// Package handlers contains the mapping to our routes
package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"github.com/Maverickme222222/api/logic"
	"github.com/Maverickme222222/api/services"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, services services.Services, log *zerolog.Logger) *mux.Router {
	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	logicI := logic.New(services)

	users := NewUserHandler(*logicI, log)
	emails := NewEmailHandler(*logicI)

	// subrouter for the /api/v1 endpoint
	s := r.PathPrefix("/api/v1").Subrouter()

	// check handlers
	// Register debug check endpoints.
	cg := NewCheckGroup(log, build)
	s.HandleFunc("/liveness", cg.Liveness).Methods(http.MethodGet)

	// subrouter for the /api/v1/users endpoint
	u := s.PathPrefix("/users").Subrouter()

	u.HandleFunc("/get", users.CreateNewUser).Methods(http.MethodGet)

	// subrouter for the /api/v1/users endpoint
	e := s.PathPrefix("/emails").Subrouter()

	e.HandleFunc("/get", emails.CreateNewEmail).Methods(http.MethodGet)

	return r
}
