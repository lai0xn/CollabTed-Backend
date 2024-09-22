package types

type WorkspaceD struct {
	Name    string `json:"workspace_name"`
	OwnerID string `json:"owner_id"`
}

type InviteUserD struct {
	Email       string `json:"email"`
	WorkspaceID string `json:"workspaceId"`
}
