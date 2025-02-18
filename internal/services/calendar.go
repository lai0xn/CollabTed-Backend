package services

import (
	"context"
	"log"
	"time"

	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
	"github.com/teambition/rrule-go"
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
		db.Event.Workspace.Link(
			db.Workspace.ID.Equals(data.WorkspaceID)),
		db.Event.Description.Set(data.Description),
		db.Event.Rrule.SetIfPresent(data.RRule),
		db.Event.Type.Set(db.EventType(data.Type)),
		db.Event.AllDay.Set(data.AllDay),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	for _, assigneeID := range data.Assignees {
		_, err = prisma.Client.UserWorkspace.FindUnique(
			db.UserWorkspace.ID.Equals(assigneeID),
		).Update(
			db.UserWorkspace.Event.Link(
				db.Event.ID.Equals(result.ID),
			),
		).Exec(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *EventService) ListEventsByWorkspace(workspaceID string, startTime, endTime time.Time) ([]types.EventInstance, error) {
	dbEvents, err := prisma.Client.Event.FindMany(
		db.Event.Workspace.Where(
			db.Workspace.ID.Equals(workspaceID),
		),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	var instances []types.EventInstance
	for _, event := range dbEvents {
		if rruleStr, ok := event.Rrule(); ok && rruleStr != "" {
			r, err := rrule.StrToRRule(rruleStr)
			if err != nil {
				log.Printf("Error parsing RRULE for event %s: %v", event.ID, err)
				continue
			}
			originalStart := event.StartTime
			r.DTStart(originalStart)

			duration := event.EndTime.Sub(originalStart)

			occurrences := r.Between(startTime, endTime, true)

			for _, occStart := range occurrences {
				occEnd := occStart.Add(duration)
				instances = append(instances, createEventInstance(event, occStart, occEnd))
			}
		} else {
			if event.StartTime.Before(endTime) && event.EndTime.After(startTime) {
				instances = append(instances, createEventInstance(event, event.StartTime, event.EndTime))
			}
		}
	}

	return instances, nil
}

func createEventInstance(event db.EventModel, start, end time.Time) types.EventInstance {
	// Handle optional description
	var description *string
	if desc, ok := event.Description(); ok && desc != "" {
		description = &desc
	}

	// Handle assignees IDs directly from the field
	assigneesIDs := event.AssineesIds

	return types.EventInstance{
		ID:           event.ID,
		Name:         event.Name,
		Description:  description,
		StartTime:    start,
		EndTime:      end,
		Type:         string(event.Type),
		CreatorID:    event.CreatorID,
		WorkspaceID:  event.WorkspaceID,
		AssigneesIds: assigneesIDs,
		MeetLink:     event.MeetLink,
	}
}
