package types

type NotifType string

const (
	MESSAGE_NOTIFICATION NotifType = "message"
	CALL_NOTIFICATION    NotifType = "call"
	KICK_NOTIFICATION    NotifType = "kick"
	JOIN_NOTIFICATION    NotifType = "join"
)

type PingNotification struct {
	Type    NotifType `json:"type"`
	Content string    `json:"content"`
	Sender  string    `json:"senderName"`
	Channel string    `json:"channelID"`
}

type CallNotification struct {
	Type     NotifType `json:"type"`
	RoomID   string    `json:"roomId"`
	CallerID string    `json:"callerId"`
}

type KickNotification struct {
	Type        NotifType `json:"type"`
	WorkspaceID string    `json:"workspaceId"`
}

type JoinUser struct {
	Type        NotifType `json:"type"`
	WorkspaceID string    `json:"workspaceId"`
}
