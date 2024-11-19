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
	Assignees   []string  `json:"assignees"`
	MeetLink    string    `json:"meetLink,omitempty"`
}
