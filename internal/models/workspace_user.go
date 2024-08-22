package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type WorkspaceUser struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"userId"`
	WorkplaceID primitive.ObjectID `bson:"workplaceId"`
	RoleID      primitive.ObjectID `bson:"roleId"`
	CreatedAt   primitive.DateTime `bson:"createdAt"`
	UpdatedAt   primitive.DateTime `bson:"updatedAt"`
}
