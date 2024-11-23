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
	startTime := data.StartTime
	endTime := data.EndTime

	result, err := prisma.Client.Event.CreateOne(
		db.Event.Name.Set(data.Name),
		db.Event.StartTime.Set(startTime),
		db.Event.EndTime.Set(endTime),
		db.Event.CreatorID.Set(data.CreatorID),
		db.Event.MeetLink.Set(data.MeetLink),
		db.Event.Description.Set(data.Description),
		db.Event.Type.Set(db.EventType(data.Type)),
		db.Event.WorkspaceID.Set(data.WorkspaceID),
		db.Event.AssineesIds.Set(data.Assignees),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	for _, assigneeID := range data.Assignees {
		usr, err := prisma.Client.UserWorkspace.FindFirst(
			db.UserWorkspace.UserID.Equals(assigneeID),
			db.UserWorkspace.WorkspaceID.Equals(data.WorkspaceID),
		).Exec(context.Background())
		if err != nil {
			return nil, err
		}
		_, err = prisma.Client.UserWorkspace.FindUnique(
			db.UserWorkspace.ID.Equals(usr.ID),
		).Update(
			db.UserWorkspace.Event.Link(
				db.Event.ID.Equals(result.ID),
			),
		).Exec(context.Background())
		if err != nil {
			return nil, err
		}

	}
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
