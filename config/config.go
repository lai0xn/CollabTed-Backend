package config

import (
	"os"

	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

var (
	JWT_SECRET string

	EMAIL_HOST     string
	EMAIL_PORT     string
	EMAIL          string
	EMAIL_PASSWORD string
)

func Load() {
	// OAuth configuration
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	types.OAuth2Configs = map[string]*types.OAuthProvider{
		"google": {
			Config: &oauth2.Config{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
				Scopes:       []string{"profile", "email"},
				Endpoint:     google.Endpoint,
			},
		},
		"facebook": {
			Config: &oauth2.Config{
				ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
				ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("FACEBOOK_REDIRECT_URL"),
				Scopes:       []string{"public_profile", "email"},
				Endpoint:     facebook.Endpoint,
			},
		},
	}

	// JWT Secret
	JWT_SECRET = os.Getenv("JWT_SECRET")

	// Email configuration
	EMAIL_HOST = os.Getenv("EMAIL_HOST")
	EMAIL_PORT = os.Getenv("EMAIL_PORT")
	EMAIL = os.Getenv("EMAIL")
	EMAIL_PASSWORD = os.Getenv("EMAIL_PASSWORD")

	// Initialize the logger
	logger.NewLogger()
}
