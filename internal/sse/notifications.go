package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/redis"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/labstack/echo/v4"
	r "github.com/redis/go-redis/v9"
)

type Notifier struct {
	client *r.Client
}

func NewNotifier() *Notifier {
	return &Notifier{
		client: redis.GetClient(),
	}
}

func (n *Notifier) NotificationHandler(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	id := c.QueryParam("userID")

	// Ensure the Content-Length is not set, as SSE is a streaming response
	c.Response().Header().Del("Content-Length")

	// Create a new Redis Pub/Sub subscriber
	pubsub := redis.GetClient().Subscribe(context.Background(), "notifs:"+id)
	defer pubsub.Close()

	_, err := pubsub.Receive(context.Background())
	if err != nil {
		log.Printf("Failed to subscribe: %v", err)
		return err
	}

	// Start a goroutine to receive messages from Redis
	ch := pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(c.Response(), "data: %s\n\n", msg.Payload)
			c.Response().Flush()
		case <-c.Request().Context().Done():
			return nil
		}
	}
}

func (n *Notifier) NotifyCallUser(userID, roomID, callerID string) error {
	call := types.CallNotification{
		Type:     types.CALL_NOTIFICATION,
		CallerID: callerID,
		RoomID:   roomID,
	}

	b, err := json.Marshal(call)
	if err != nil {
		log.Printf("Failed to marshal call: %v", err)
		return err
	}
	err = n.client.Publish(context.Background(), "notifs:"+userID, b).Err()
	if err != nil {
		log.Printf("Failed to publish notification: %v", err)
		return err
	}
	return nil
}

func (n *Notifier) NotifyPing(userID string, notif types.PingNotification) error {
	b, err := json.Marshal(notif)
	if err != nil {
		log.Printf("Failed to marshal call :%v", err)
		return err
	}
	fmt.Println(notif, "notification")
	err = n.client.Publish(context.Background(), "notifs:"+userID, b).Err()
	if err != nil {
		log.Printf("Failed to publish notification: %v", err)
		return err
	}
	return nil
}

func (n *Notifier) NotifyKickUser(userID, workspaceID string) error {
	notif := types.KickNotification{
		Type:        types.KICK_NOTIFICATION,
		WorkspaceID: workspaceID,
	}

	b, err := json.Marshal(notif)
	if err != nil {
		log.Printf("Failed to marshal call: %v", err)
		return err
	}

	logger.LogInfo().Msgf("Publishing kick notification: %s", string(b))
	err = n.client.Publish(context.Background(), "notifs:"+userID, b).Err()
	if err != nil {
		log.Printf("Failed to publish notification: %v", err)
		return err
	}
	return nil
}

func (n *Notifier) NotifyJoinUser(userID, workspaceID string) error {
	notif := types.JoinUser{
		Type:        types.JOIN_NOTIFICATION,
		WorkspaceID: workspaceID,
	}

	b, err := json.Marshal(notif)
	if err != nil {
		log.Printf("Failed to marshal call: %v", err)
		return err
	}

	logger.LogInfo().Msgf("Publishing join notification: %s", string(b))
	err = n.client.Publish(context.Background(), "notifs:"+userID, b).Err()
	if err != nil {
		log.Printf("Failed to publish notification: %v", err)
		return err
	}
	return nil
}