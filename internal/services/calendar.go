package services

import (
	"context"

	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
)

type EventService struct{}

func NewEventService() *EventService {
	return &EventService{}
}

func (s *EventService) CreateEvent(data types.EventD) (*db.EventModel, error) {
	result, err := prisma.Client.Event.CreateOne(
		db.Event.Name.Set(data.Name),
		db.Event.StartTime.Set(data.StartTime),
		db.Event.EndTime.Set(data.EndTime),
		db.Event.Workspace.Link(
			db.Workspace.ID.Equals(data.WorkspaceID),
		),
		db.Event.Description.Set(data.Description),
		db.Event.Type.Set(db.EventType(data.Type)),
	).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *EventService) ListEventsByWorkspace(workspaceID string) ([]db.EventModel, error) {
	result, err := prisma.Client.Event.FindMany(
		db.Event.Workspace.Where(
			db.Workspace.ID.Equals(workspaceID),
		),
	).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return result, nil
}
