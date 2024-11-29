package types

type PingNotification struct {
	Content 	string  `json:"content"`
	SenderID 	string  `json:"senderId"`
	ChannelID	string	`json:"channelId"`
}

type CallNotification struct {
	RoomID   string    `json:"roomId"`
	CallerID string	   `json:"callerId"`
}
