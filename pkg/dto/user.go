package dto

type UserD struct {
	Name        string `validate:"required" json:"name"`
	Email       string `validate:"required;email" json:"email"`
	PhoneNumber string `validate:"required" json:"number"`
	Password    string `valdiate:"required" json:"password"`
}
