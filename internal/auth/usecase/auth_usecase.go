package usecase

import (
	"context"
	"time"

	"github.com/CollabTed/CollabTed-Backend/internal/auth/repository"
	"github.com/CollabTed/CollabTed-Backend/internal/models"
	"github.com/CollabTed/CollabTed-Backend/pkg/consts"
	"github.com/CollabTed/CollabTed-Backend/pkg/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(ctx context.Context, user dto.UserD) error
	Login(ctx context.Context, email, password string) (*models.User, error)
}

type authUseCase struct {
	repo repository.AuthRepository
}

func NewAuthUseCase(repo repository.AuthRepository) AuthUseCase {
	return &authUseCase{repo: repo}
}

func (u *authUseCase) Register(ctx context.Context, user dto.UserD) error {
	existingUser, _ := u.repo.GetUserByEmail(ctx, user.Email)
	if existingUser != nil {
		return consts.EMAIL_IN_USE
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userModel := models.User{
		ID:          primitive.NewObjectID(),
		Email:       user.Email,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Password:    string(hashedPassword),
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}
	return u.repo.CreateUser(ctx, &userModel)
}

func (u *authUseCase) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, consts.INVALID_CREDENTIALS
	}

	return user, nil
}
