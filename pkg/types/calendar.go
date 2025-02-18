package types

import "time"

type EventD struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"startTime" validate:"required"`
	EndTime     time.Time `json:"endTime" validate:"required"`
	CreatorID   string    `json:"creatorId" validate:"required"`
	Type        string    `json:"type" validate:"required,oneof=EVENT MEET WORKING_HOURS"`
	WorkspaceID string    `json:"workspaceId" validate:"required"`
	Assignees   []string  `json:"assineesIds"`
	MeetLink    string    `json:"meetLink,omitempty"`
	RRule       *string   `json:"rrule,omitempty"`
	AllDay      bool      `json:"allDay"`
}

type EventInstance struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	Type         string    `json:"type"`
	CreatorID    string    `json:"creatorId"`
	WorkspaceID  string    `json:"workspaceId"`
	AssigneesIds []string  `json:"assigneesIds"`
	MeetLink     string    `json:"meetLink"`
}
