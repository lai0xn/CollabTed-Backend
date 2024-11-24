package types

import (
	"encoding/json"
	"time"
)

// TaskD is the data structure used for creating or updating a task.
type TaskD struct {
	Title        string            `json:"title"`        // Task title
	Description  []json.RawMessage `json:"description"`  // Task description
	DueDate      time.Time         `json:"dueDate"`      // Task due date
	Priority     string            `json:"priority"`     // Priority (e.g., HIGH, MEDIUM, LOW)
	StatusID     string            `json:"statusId"`     // Status ID for task status
	ProjectID    string            `json:"projectId"`    // Project ID the task belongs to
	AssigneesIDs []string          `json:"assigneesIds"` // List of user IDs assigned to the task
	WorkspaceID  string            `json:"workspaceId"`  // Workspace ID to check user permissions
}
