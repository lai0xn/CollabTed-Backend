package types

import "golang.org/x/oauth2"

type OAuthProvider struct {
	Config   *oauth2.Config
	Endpoint oauth2.Endpoint
}

var OAuth2Configs map[string]*OAuthProvider

type OAuthUser struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"picture"`
}
