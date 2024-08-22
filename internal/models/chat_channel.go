package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChatChannel struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	ChannelName  string               `bson:"channelName"`
	WorkplaceID  primitive.ObjectID   `bson:"workplaceId"`
	Participants []primitive.ObjectID `bson:"participants,omitempty"`
	Attachments  []primitive.ObjectID `bson:"attachments,omitempty"`
	Messages     []primitive.ObjectID `bson:"messages,omitempty"`
	CreatedAt    primitive.DateTime   `bson:"createdAt"`
	UpdatedAt    primitive.DateTime   `bson:"updatedAt"`
}
