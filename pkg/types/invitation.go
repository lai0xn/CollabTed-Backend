package types

type InvitationStatus string

const (
	PENDING  InvitationStatus = "PENDING"
	ACCEPTED InvitationStatus = "ACCEPTED"
	DECLINED InvitationStatus = "DECLINED"
)

type InviteUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	WorkspaceID string `json:"workspaceId" validate:"required"`
}

type AcceptInviteRequest struct {
	Token string `json:"token" validate:"required"`
}

type InvitationD struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	Status      string `json:"status"`
	WorkspaceID string `json:"workspaceId"`
}
