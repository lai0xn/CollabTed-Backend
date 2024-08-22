package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workspace struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	WorkspaceName string               `bson:"workplaceName"`
	CreatedAt     primitive.DateTime   `bson:"createdAt"`
	UpdatedAt     primitive.DateTime   `bson:"updatedAt"`
	Channels      []primitive.ObjectID `bson:"channels,omitempty"`
	VideoCalls    []primitive.ObjectID `bson:"videoCalls,omitempty"`
	NewField      primitive.ObjectID   `bson:"newField,omitempty"`
}
