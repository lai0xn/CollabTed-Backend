package types

type ChannelD struct {
	Name         string   `json:"name"`
	WorkspaceID  string   `json:"workspaceID"`
	Participants []string `json:"participants"`
}

type MessageD struct {
	SenderID  string `json:"senderID"`
	ChannelID string `json:"channelID"`
	Content   string `json:"content"`
}
