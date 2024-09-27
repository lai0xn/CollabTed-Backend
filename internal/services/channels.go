package services

import (
	"context"
	"fmt"

	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
)

type ChannelService struct{}

func NewChannelService() *ChannelService {
	return &ChannelService{}
}

// CreateChannel creates a new channel in a workspace and adds existing participants.
func (s *ChannelService) CreateChannel(data types.ChannelD) (*db.ChannelModel, error) {
	// Create a new channel
	result, err := prisma.Client.Channel.CreateOne(
		db.Channel.Name.Set(data.Name),
		db.Channel.Workspace.Link(
			db.Workspace.ID.Equals(data.WorkspaceID),
		),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	// Link participants to the channel and workspace
	for _, participantID := range data.Participants {
		// Ensure that the user exists in the workspace and link them to the channel
		_, err := prisma.Client.UserWorkspace.FindUnique(
			db.UserWorkspace.ID.Equals(participantID),
		).Update(
			db.UserWorkspace.Channel.Link(
				db.Channel.ID.Equals(result.ID), // Link the created channel
			),
		).Exec(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to add participant with ID %s to the channel: %v", participantID, err)
		}
	}

	return result, nil
}

// GetChannelById retrieves a channel by its ID.
func (s *ChannelService) GetChannelById(channelID string) (*db.ChannelModel, error) {
	channel, err := prisma.Client.Channel.FindUnique(
		db.Channel.ID.Equals(channelID),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return channel, nil
}

// ListChannelsByWorkspace lists all channels in a workspace.
func (s *ChannelService) ListChannelsByWorkspace(workspaceID string) ([]db.ChannelModel, error) {
	channels, err := prisma.Client.Channel.FindMany(
		db.Channel.WorkspaceID.Equals(workspaceID),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *ChannelService) AddParticipant(workspaceID, channelID, userID string) (*db.UserWorkspaceModel, error) {
	ctx := context.Background()
	user, err := prisma.Client.UserWorkspace.FindFirst(
		db.UserWorkspace.UserID.Equals(userID),
		db.UserWorkspace.WorkspaceID.Equals(workspaceID),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	_, err = prisma.Client.Channel.FindUnique(
		db.Channel.ID.Equals(channelID),
	).Update(
		db.Channel.Participants.Link(db.UserWorkspace.ID.Equals(user.ID)),
	).Exec(ctx)
	return user, nil
}
