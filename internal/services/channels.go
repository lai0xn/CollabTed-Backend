package services

import (
	"context"
	"encoding/json"

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

	return result, nil
}

// GetChannelById retrieves a channel by its ID.
func (s *ChannelService) GetChannelById(channelID string) (*db.ChannelModel, error) {
	channel, err := prisma.Client.Channel.FindUnique(
		db.Channel.ID.Equals(channelID),
	).With(db.Channel.Participants.Fetch()).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return channel, nil
}

// ListChannelsByWorkspace lists all channels in a workspace.
func (s *ChannelService) ListChannelsByWorkspace(workspaceID string) ([]db.ChannelModel, error) {
	channels, err := prisma.Client.Channel.FindMany(
		db.Channel.WorkspaceID.Equals(workspaceID),
	).With(db.Channel.Participants.Fetch()).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *ChannelService) AddParticipants(workspaceID string, channelID string, userIDs []json.RawMessage) ([]*db.UserWorkspaceModel, error) {
	ctx := context.Background()
	var addedUsers []*db.UserWorkspaceModel

	for _, userID := range userIDs {
		// Find the user in the workspace
		user, err := prisma.Client.UserWorkspace.FindFirst(
			db.UserWorkspace.UserID.Equals(string(userID)),
			db.UserWorkspace.WorkspaceID.Equals(workspaceID),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}

		// Add the user to the channel's participants
		_, err = prisma.Client.Channel.FindUnique(
			db.Channel.ID.Equals(channelID),
		).Update(
			db.Channel.Participants.Link(db.UserWorkspace.ID.Equals(user.ID)),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}

		// Append the added user to the result list
		addedUsers = append(addedUsers, user)
	}

	return addedUsers, nil
}
