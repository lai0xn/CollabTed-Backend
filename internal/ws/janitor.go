package ws

import (
	"fmt"
	"log"
	"sync"

	"github.com/CollabTED/CollabTed-Backend/prisma/db"
	"github.com/gorilla/websocket"
)

type Connection struct {
	conn        *websocket.Conn
	name        string
	workspaceID string
	userID      string
}

type Message struct {
	SenderID  string `json:"senderID"`
	ChannelID string `json:"channelID"`
	Content   string `json:"content"`
	Recievers []db.UserWorkspaceModel
}

var (
	connection = make(chan Connection)
	messages   = make(chan Message)
	closing    = make(chan string)
	users      = make(map[string]Connection)
	mu         sync.RWMutex
)

func Hub() {
	for {
		select {
		case con := <-connection:
			fmt.Println(fmt.Sprintf("user: %s connected", con.userID))
			mu.Lock()
			users[con.userID] = con
			mu.Unlock()
		case msg := <-messages:
			for _, user := range msg.Recievers {
				con, ok := users[user.UserID]
				if !ok {
					continue
				}
				err := con.conn.WriteJSON(msg)
				if err != nil {
					log.Println(err.Error())
				}
			}
		case id := <-closing:
			fmt.Println(fmt.Sprintf("user: %s disconnected", id))
			mu.Lock()
			delete(users, id)
			mu.Unlock()

		}
	}
}
