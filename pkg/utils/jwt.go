package utils

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CollabTED/CollabTed-Backend/config"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(id string, email string, name string, profilePicture string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.Claims{
		ID:             id,
		Name:           name,
		Email:          email,
		ProfilePicture: profilePicture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	})
	fmt.Println(config.JWT_SECRET)
	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetJWTAsCookie(w http.ResponseWriter, id string, email string, name string, profilePicture string) error {
	jwtToken, err := GenerateJWT(id, email, name, profilePicture)
	if err != nil {
		logger.LogError().Msgf("Error generating JWT: %v", err)
		return err
	}

	if jwtToken == "" {
		errMsg := "Generated JWT token is empty"
		log.Println(errMsg)
		return http.ErrNoCookie
	}

	cookie := &http.Cookie{
		Name:     "jwt",                          // Cookie name
		Value:    jwtToken,                       // JWT token
		Expires:  time.Now().Add(72 * time.Hour), // Same as JWT expiration
		HttpOnly: true,                           // Ensure cookie is HttpOnly
		Secure:   config.SECURE_COOKIE,           // Set to true in production (requires HTTPS)
		Path:     "/",                            // Cookie path
	}

	http.SetCookie(w, cookie)
	return nil
}

func GenerateInvitationToken() (string, error) {
	toekn := uuid.NewString()
	return toekn, nil
}

func DeleteJWTCookie(w http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   config.SECURE_COOKIE,
	}
	http.SetCookie(w, cookie)
	return nil
}
