package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type VideoCall struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	CallType      []string             `bson:"callType"`
	Participants  []primitive.ObjectID `bson:"participants,omitempty"`
	ScheduledTime primitive.DateTime   `bson:"scheduledTime"`
	CallLink      string               `bson:"callLink"`
	CreatedBy     primitive.ObjectID   `bson:"createdBy"`
	Status        string               `bson:"status"`
	Messages      []primitive.ObjectID `bson:"messages,omitempty"`
	CreatedAt     primitive.DateTime   `bson:"createdAt"`
	UpdatedAt     primitive.DateTime   `bson:"updatedAt"`
}
