package types

type AttachmentD struct {
	ChannelID   string `json:"channelID"`
	WorkspaceID string `json:"workspaceID"`
	SenderID    string `json:"senderID"`
	File        string `json:"file"`
	Title       string `json:"title"`
}
