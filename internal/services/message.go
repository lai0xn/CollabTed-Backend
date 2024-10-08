package services

import (
	"context"
	"fmt"
	"time"

	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) SendMessage(data types.MessageD) (*db.MessageModel, error) {
	message, err := prisma.Client.Message.CreateOne(
		db.Message.Content.Set(data.Content),
		db.Message.Channel.Link(
			db.Channel.ID.Equals(data.ChannelID),
		),
		db.Message.Sender.Link(
			db.User.ID.Equals(data.SenderID),
		),
		db.Message.CreatedAt.Set(time.Now()),
	).Exec(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %v", err)
	}

	return message, nil
}

func (s *MessageService) GetMessagesByChannel(channelID string) ([]db.MessageModel, error) {
	messages, err := prisma.Client.Message.FindMany(
		db.Message.ChannelID.Equals(channelID),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *MessageService) GetMessageById(messageID string) (*db.MessageModel, error) {
	message, err := prisma.Client.Message.FindUnique(
		db.Message.ID.Equals(messageID),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageService) DeleteMessage(messageID string) error {
	_, err := prisma.Client.Message.FindUnique(
		db.Message.ID.Equals(messageID),
	).Delete().Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
