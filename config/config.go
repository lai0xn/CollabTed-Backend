package config

import (
	"log"
	"os"
	"strconv"

	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
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
	MONGO_URI          string
)

func Load() {
	// Initialize the logger
	logger.NewLogger()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Logger.Err(err).Msg("Error loading .env file, using system environment variables instead.")
	}

	// OAuth2 Configurations
	initOAuthConfigs()

	// Load environment-specific settings
	JWT_SECRET = mustGetEnv("JWT_SECRET")

	EMAIL_HOST = mustGetEnv("EMAIL_HOST")
	EMAIL_PORT = mustGetEnv("EMAIL_PORT")
	EMAIL = mustGetEnv("EMAIL")
	EMAIL_PASSWORD = mustGetEnv("EMAIL_PASSWORD")

	SECURE_COOKIE = mustParseBool("SECURE_COOKIE", true)
	HOST_URL = mustGetEnv("HOST_URL")
	ALLOWED_ORIGINS = mustGetEnv("ALLOWED_ORIGINS")

	LIVEKIT_API_KEY = mustGetEnv("LIVEKIT_API_KEY")
	LIVEKIT_API_SECRET = mustGetEnv("LIVEKIT_API_SECRET")
	CLOUDINARY_URL = mustGetEnv("CLOUDINARY_URL")

	MONGO_URI = getMongoURI()
	logger.Logger.Info().Msgf("Starting %s environment", os.Getenv("APP_ENV"))
}

func initOAuthConfigs() {
	types.OAuth2Configs = map[string]*types.OAuthProvider{
		"google": {
			Config: &oauth2.Config{
				ClientID:     mustGetEnv("GOOGLE_CLIENT_ID"),
				ClientSecret: mustGetEnv("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  mustGetEnv("GOOGLE_REDIRECT_URL"),
				Scopes:       []string{"profile", "email"},
				Endpoint:     google.Endpoint,
			},
		},
	}
}

func getMongoURI() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		logger.Logger.Warn().Msg("APP_ENV not set, defaulting to 'dev'")
		env = "dev"
	}

	var mongoURI string
	if env == "prod" {
		mongoURI = os.Getenv("MONGO_URI_PROD")
	} else {
		mongoURI = os.Getenv("MONGO_URI_DEV")
	}

	if mongoURI == "" {
		log.Fatalf("MongoDB URI for %s environment is not set", env)
	}

	os.Setenv("MONGO_URI", mongoURI)

	return mongoURI
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func mustParseBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		logger.Logger.Warn().Msgf("Invalid boolean value for %s: %s, defaulting to %t", key, value, defaultValue)
		return defaultValue
	}
	return parsed
}
