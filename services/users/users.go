// Package users contains the functionality to access users services.
package users

import (
	"context"
	"fmt"

	"github.com/Maverickme222222/users/usermgmt"
)

// Service is the access point for the remote configuration service.
type Service struct {
	client usermgmt.UserManagementClient
}

// New returns an instance of the access point to the remote configuration service via GRPC.
func New(client usermgmt.UserManagementClient) *Service {
	return &Service{
		client: client,
	}
}

// CreateNewUser  implements the Service interface for CreateNewUser.
func (s *Service) CreateNewUser(ctx context.Context, name string) (*usermgmt.NewUserResponse, error) {

	req := &usermgmt.NewUser{
		Name: name,
	}

	res, err := s.client.CreateNewUser(ctx, req)
	if err != nil {
		fmt.Printf("ERRRS %+v", err)
		return nil, err
	}

	return res, nil
}
