package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type User struct {
	UserID      string
	WorkspaceID string
	Conn        *websocket.Conn
}

type OnlineEvent struct {
	UserID string `json:"userID"`
	Event  string `json:"event"`
}

var (
	online        = make(chan User)
	disconnected  = make(chan User)
	workspaceLock sync.Mutex
	workspaces    = make(map[string][]User)
)

// BroadcastEvent sends an event to all users in a workspace.
func broadcastEvent(workspaceID string, event OnlineEvent) {
	workspaceLock.Lock()
	defer workspaceLock.Unlock()

	users, ok := workspaces[workspaceID]
	if !ok {
		return
	}

	for _, user := range users {
		if err := user.Conn.WriteJSON(event); err != nil {
			log.Printf("Failed to send event to user %s: %v", user.UserID, err)
		}
	}
}

// WatchConnect handles new user connections.
func WatchConnect() {
	for user := range online {
		workspaceLock.Lock()

		// Add the user to the workspace
		workspaces[user.WorkspaceID] = append(workspaces[user.WorkspaceID], user)

		workspaceLock.Unlock()

		// Broadcast the "connected" event to all users in the workspace
		broadcastEvent(user.WorkspaceID, OnlineEvent{
			UserID: user.UserID,
			Event:  "connected",
		})
	}
}

func WatchDisconnect() {
	for user := range disconnected {
		workspaceLock.Lock()

		// Remove the user from the workspace
		users, ok := workspaces[user.WorkspaceID]
		if ok {

			updatedUsers := make([]User, 0, len(users)-1)
			for _, u := range users {
				if u.UserID != user.UserID {
					updatedUsers = append(updatedUsers, u)
				}
			}
			workspaces[user.WorkspaceID] = updatedUsers
		}

		workspaceLock.Unlock()

		// Broadcast the "disconnected" event to all users in the workspace
		broadcastEvent(user.WorkspaceID, OnlineEvent{
			UserID: user.UserID,
			Event:  "disconnected",
		})
	}
}
