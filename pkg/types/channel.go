package types

type ChannelD struct {
	Name         string   `json:"name"`
	WorkspaceID  string   `json:"workspaceID"`
	Participants []string `json:"participants"`
}

type ParticipantD struct {
	WorkspaceID string `json:"workspaceID"`
	UserID      string `json:"userID"`
	ChannelD    string `json:"channelID"`
}

type MessageD struct {
	SenderID  string `json:"senderID"`
	ChannelID string `json:"channelID"`
	Content   string `json:"content"`
}
