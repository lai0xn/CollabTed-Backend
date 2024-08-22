package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Content     string               `bson:"content"`
	SenderID    primitive.ObjectID   `bson:"senderId"`
	Attachments []primitive.ObjectID `bson:"attachments,omitempty"`
}
