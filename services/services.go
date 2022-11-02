// Package services centralizes all the downstream service wrappers to make it easier to initialize across various parts of the kappa stack.
package services

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/Maverickme222222/api/services/emails"
	"github.com/Maverickme222222/api/services/users"
	emailsGRPC "github.com/Maverickme222222/emails/emailmgmt"
	usersGRPC "github.com/Maverickme222222/users/usermgmt"
)

type (
	// Services is an aggregation of all the downstream services
	Services struct {
		Emails *emails.Service
		Users  *users.Service
	}
)

// Register uses the provided configs to configure service dependencies.
func Register(
	ctx context.Context,
	usersConf string,
	emailsConf string) (Services, error) {

	var services = Services{}

	conn, err := grpc.Dial(emailsConf, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect to email")
	}

	services.Emails = emails.New(emailsGRPC.NewEmailManagementClient(conn))

	// usersConn Service
	usersConn, err := grpc.Dial(usersConf, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect to users")
	}
	services.Users = users.New(usersGRPC.NewUserManagementClient(usersConn))

	return services, nil
}
