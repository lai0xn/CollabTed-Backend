package auth

import (
	"context"

	"github.com/CollabTed/CollabTed-Backend/internal/models"
	"github.com/CollabTed/CollabTed-Backend/pkg/dto"
)

type AuthUseCase interface {
	Register(ctx context.Context, user dto.UserD) error
	Login(ctx context.Context, email, password string) (*models.User, error)
}
