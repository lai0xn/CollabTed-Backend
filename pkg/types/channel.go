package types

type ChannelD struct {
	Name            string   `json:"name"`
	WorkspaceID     string   `json:"workspaceID"`
	CreatorID       string   `json:"creatorID"`
	ParticipantsIDS []string `json:"participantsIds"`
}

type ParticipantD struct {
	WorkspaceID string   `json:"workspaceID"`
	UsersID     []string `json:"usersID"`
	ChannelID   string   `json:"channelID"`
}

type MessageD struct {
	SenderID  string `json:"senderID"`
	ChannelID string `json:"channelID"`
	Content   string `json:"content"`

	IsReply         bool   `json:"isReply"`
	ReplyToMessage  string `json:"replyToMessage"`
	ReplyToUserName string `json:"replyToUserName"`

	AttachmentLink  string `json:"attachmentLink"`
	AttachmentTitle string `json:"attachmentTitle"`
}
