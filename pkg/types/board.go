package types

import (
	"encoding/json"
)

type BoardD struct {
	Elements    []json.RawMessage `json:"elements"`
	WorkspaceID string            `json:"workspaceId"`
}
