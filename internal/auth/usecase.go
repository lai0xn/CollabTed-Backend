package auth

import (
	"context"
	"errors"
	"time"

	"github.com/CollabTed/CollabTed-Backend/internal/auth/repository"
	"github.com/CollabTed/CollabTed-Backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (*models.User, error)
}

type authUseCase struct {
	repo repository.AuthRepository
}

func NewAuthUseCase(repo repository.AuthRepository) AuthUseCase {
	return &authUseCase{repo: repo}
}

func (u *authUseCase) Register(ctx context.Context, user *models.User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("name, email, and password are required")
	}

	existingUser, _ := u.repo.GetUserByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Set default values
	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return u.repo.CreateUser(ctx, user)
}

func (u *authUseCase) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
