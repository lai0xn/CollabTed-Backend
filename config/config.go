package config

import (
	"log"
	"os"
	"strconv"

	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

var (
	JWT_SECRET      string
	EMAIL_HOST      string
	EMAIL_PORT      string
	EMAIL           string
	EMAIL_PASSWORD  string
	SECURE_COOKIE   bool
	HOST_URL        string
	ALLOWED_ORIGINS string
)

func Load() {
	// Initialize the logger
	logger.NewLogger()

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

	// Secure Cookie
	SECURE_COOKIE, err = strconv.ParseBool(os.Getenv("SECURE_COOKIE"))
	if err != nil {
		log.Printf("Error parsing SECURE_COOKIE environment variable: %v", err)
		SECURE_COOKIE = false
	}

	// HOST URL
	HOST_URL = os.Getenv("HOST_URL")

	// Allowed Origins
	ALLOWED_ORIGINS = os.Getenv("ALLOWED_ORIGINS")
}
