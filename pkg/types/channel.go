package types

type ChannelD struct {
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceID"`
}

type ParticipantD struct {
	WorkspaceID string `json:"workspaceID"`
	UserID      string `json:"userID"`
	ChannelID   string `json:"channelID"`
}

type MessageD struct {
	SenderID  string `json:"senderID"`
	ChannelID string `json:"channelID"`
	Content   string `json:"content"`
}
