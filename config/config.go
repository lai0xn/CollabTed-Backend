package config

import (
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
	JWT_SECRET         string
	EMAIL_HOST         string
	EMAIL_PORT         string
	EMAIL              string
	EMAIL_PASSWORD     string
	SECURE_COOKIE      bool
	HOST_URL           string
	ALLOWED_ORIGINS    string
	LIVEKIT_API_KEY    string
	LIVEKIT_API_SECRET string
	CLOUDINARY_URL     string
)

func Load() {
	// Initialize the logger
	logger.NewLogger()

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Err(err).Msg("Error loading .env file")
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
		SECURE_COOKIE = true
	}
	logger.Logger.Info().Msg("Secure Cookie: " + strconv.FormatBool(SECURE_COOKIE))

	// HOST URL
	HOST_URL = os.Getenv("HOST_URL")
	logger.Logger.Info().Msg(HOST_URL)

	// Allowed Origins
	ALLOWED_ORIGINS = os.Getenv("ALLOWED_ORIGINS")
	logger.Logger.Info().Msg(ALLOWED_ORIGINS)

	// Live Kit Credentials
	LIVEKIT_API_KEY = os.Getenv("LIVEKIT_API_KEY")
	LIVEKIT_API_SECRET = os.Getenv("LIVEKIT_API_SECRET")

	// Cloudinary URL
	CLOUDINARY_URL = os.Getenv("CLOUDINARY_URL")
}
