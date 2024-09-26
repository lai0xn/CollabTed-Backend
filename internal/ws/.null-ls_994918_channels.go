package ws

import (
	"fmt"
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

type WsChatHandler struct {
	srv services.ChannelService
}

func (ws WsChatHandler) Chat(c echo.Context) error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
	}
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*types.Claims)
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer conn.Close()
	connection <- Connection{
		conn: conn,
		name: "TestUser",
	}
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("Connection Closed")
		return nil
	})
	for {
		var data map[string]interface{}
		err := conn.ReadJSON(&data)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = conn.WriteJSON(data)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
}
