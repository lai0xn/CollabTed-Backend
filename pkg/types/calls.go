package types

type Call struct {
	RoomID   string `json:"roomId"`
	CallerID string `json:"callerId"`
}
