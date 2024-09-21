package types

type RegisterPayload struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required"`
	Password       string `json:"password" validate:"required"`
	ProfilePicture string `json:"profile_picture"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
