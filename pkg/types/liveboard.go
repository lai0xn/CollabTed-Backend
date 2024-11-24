package types

type LiveBoardD struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Users       []string `json:"users"`
	WorkspaceId string   `json:"workspaceId"`
}
