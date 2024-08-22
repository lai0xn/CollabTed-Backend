package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	Email          string             `bson:"email"`
	PhoneNumber    string             `bson:"phoneNumber,omitempty"`
	Password       string             `bson:"password"`
	ProfilePicture string             `bson:"profilePicture,omitempty"`
	CreatedAt      primitive.DateTime `bson:"createdAt"`
	UpdatedAt      primitive.DateTime `bson:"updatedAt"`
}
