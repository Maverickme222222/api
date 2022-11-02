package logic

import (
	"context"

	"github.com/Maverickme222222/api/services"
)

// Logic defines the business logic for relay related operations within this service
type Logic struct {
	services services.Services
}

// New creates a new instance of Logic
func New(services services.Services) *Logic {
	return &Logic{
		services: services,
	}
}

func (l *Logic) CreateNewUser(ctx context.Context, name string) (string, error) {
	res, _ := l.services.Users.CreateNewUser(ctx, name)
	return res.GetName(), nil
}

func (l *Logic) CreateNewEmail(ctx context.Context, name string) (string, error) {
	res, _ := l.services.Emails.CreateNewEmail(ctx, name)
	return res.GetName(), nil
}
