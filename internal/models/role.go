package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	RoleName    string               `bson:"roleName"`
	Permissions []primitive.ObjectID `bson:"permissions,omitempty"`
	CreatedAt   primitive.DateTime   `bson:"createdAt"`
	UpdatedAt   primitive.DateTime   `bson:"updatedAt"`
}
