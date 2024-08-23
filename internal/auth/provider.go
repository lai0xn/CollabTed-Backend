package auth

import (
	"github.com/CollabTed/CollabTed-Backend/internal/auth/http"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeAuthHandler(db *mongo.Database) *http.AuthHandler {
	authRepository := NewMongoAuthRepository(db)
	usecaseAuthUseCase := NewAuthUseCase(authRepository)
	authHandler := http.NewAuthHandler(usecaseAuthUseCase)
	return authHandler
}
