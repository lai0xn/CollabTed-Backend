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
	WriteMu     sync.Mutex
}

type OnlineEvent struct {
	UserID string   `json:"userID"`
	Event  string   `json:"event"`
	Users  []string `json:"users,omitempty"`
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
		go func(u User) {
			u.WriteMu.Lock()
			defer u.WriteMu.Unlock()

			if err := u.Conn.WriteJSON(event); err != nil {
				log.Printf("Failed to send event to user %s: %v", u.UserID, err)
			}
		}(user)
	}
}

// WatchConnect handles new user connections.
func WatchConnect() {
	for user := range online {
		workspaceLock.Lock()

		// Get existing users BEFORE adding new user
		existingUsers := workspaces[user.WorkspaceID]

		// Add the new user to the workspace
		workspaces[user.WorkspaceID] = append(existingUsers, user)

		workspaceLock.Unlock()

		// Send initial users to new connection
		if len(existingUsers) > 0 {
			initialUserIDs := make([]string, len(existingUsers))
			for i, u := range existingUsers {
				initialUserIDs[i] = u.UserID
			}

			initialEvent := OnlineEvent{
				Event: "initial",
				Users: initialUserIDs,
			}

			user.WriteMu.Lock()
			err := user.Conn.WriteJSON(initialEvent)
			user.WriteMu.Unlock()

			if err != nil {
				log.Printf("Failed to send initial users to %s: %v", user.UserID, err)
			}
		}

		// Broadcast the "connected" event to everyone
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
