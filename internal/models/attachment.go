package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Attachment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FileName   string             `bson:"fileName"`
	FileUrl    string             `bson:"fileUrl"`
	UploaderID primitive.ObjectID `bson:"uploaderId"`
}
