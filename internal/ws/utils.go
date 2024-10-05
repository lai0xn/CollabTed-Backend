package ws

import "github.com/CollabTED/CollabTed-Backend/prisma/db"

// These functions can be used from any part of the code to send data to the users real time

// SendNotification sends a notification message to the specified recipients.
func SendNotification(senderID string, channelID string, content string, recipients []db.UserWorkspaceModel) {
	msg := Message{
		SenderID:  senderID,
		ChannelID: channelID,
		Content:   content,
		Type:      MessageTypeNotification,
		Recievers: recipients,
	}
	// Push the notification message into the messages channel for broadcasting
	messages <- msg
}

// SendMessage sends a chat message to the specified recipients.
func SendMessage(senderID string, channelID string, content string, recipients []db.UserWorkspaceModel) {
	msg := Message{
		SenderID:  senderID,
		ChannelID: channelID,
		Content:   content,
		Type:      MessageTypeBroadcast, // Use a broadcast message type or private, depending on your use case
		Recievers: recipients,
	}
	// Push the message into the messages channel for broadcasting
	messages <- msg
}
