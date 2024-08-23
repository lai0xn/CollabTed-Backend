package repository

import (
	"context"

	"github.com/CollabTed/CollabTed-Backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
}
