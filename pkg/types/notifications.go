package types

type NotifType string

const (
	MESSAGE_NOTIFICATION NotifType = "message"
	CALL_NOTIFICATION    NotifType = "call"
)

type PingNotification struct {
	Type    NotifType `json:"type"`
	Content string    `json:"content"`
	Sender  string    `json:"senderName"`
	Channel string    `json:"channelName"`
}

type CallNotification struct {
	Type     NotifType `json:"type"`
	RoomID   string    `json:"roomId"`
	CallerID string    `json:"callerId"`
}
