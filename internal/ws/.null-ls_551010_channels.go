package ws

import (
	"fmt"
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/config"
	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
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
	token_string := c.QueryParam("token")

	token := jwt.ParseWithClaims(token_string, types.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET), nil
	})
	claims := token.Claims.(*types.Claims)
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer conn.Close()
	connection <- Connection{
		conn:   conn,
		name:   claims.Name,
		userID: claims.ID,
	}
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("disconnected")
		closing <- claims.ID
		return nil
	})
	for {
		var data Message
		err := conn.ReadJSON(&data)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		channel, err := ws.srv.GetChannelById(data.ChannelID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		data.Recievers = channel.Participants()
		messages <- data
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
}
