package usecase

import (
	"context"

	"github.com/CollabTed/CollabTed-Backend/internal/models"
)

type AuthUseCase interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (*models.User, error)
}
