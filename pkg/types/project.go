package types

type ProjectD struct {
	Title        string   `json:"title"`
	WorksapceID  string   `json:"workspaceID"`
	LeadID       string   `json:"leadID"`
	AssigneesIDs []string `json:"assigneesIDs"`
}
