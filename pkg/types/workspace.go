package types

type WorkspaceD struct {
	Name    string `json:"workspace_name"`
	OwnerID string `json:"owner_id"`
}

type InviteUserD struct {
	Email       string `json:"email"`
	WorkspaceID string `json:"workspaceId"`
}

type UserWorkspace struct {
	UserWorkspaceID string `json:"userWorkspaceID"`
	ID              string `json:"id"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	ProfilePicture  string `json:"profilePicture"`
	Role            string `json:"role"`
}
