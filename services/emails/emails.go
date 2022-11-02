// Package emails contains the functionality to access users services.
package emails

import (
	"context"

	"github.com/Maverickme222222/emails/emailmgmt"
)

// Service is the access point for the remote configuration service.
type Service struct {
	client emailmgmt.EmailManagementClient
}

// New returns an instance of the access point to the remote configuration service via GRPC.
func New(client emailmgmt.EmailManagementClient) *Service {
	return &Service{
		client: client,
	}
}

// CreateNewEmail  implements the Service interface for CreateNewUser.
func (s *Service) CreateNewEmail(ctx context.Context, name string) (*emailmgmt.NewEmailResponse, error) {

	req := &emailmgmt.NewEmail{
		Name: name,
	}

	res, err := s.client.CreateNewEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
