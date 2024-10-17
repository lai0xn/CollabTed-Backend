package ws

import (
	"fmt"
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/config"
	"github.com/CollabTED/CollabTed-Backend/internal/services"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
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
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	cookie, err := c.Cookie("jwt")
	fmt.Println("cookie", cookie)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if cookie.Value == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "token is required")
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
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
		logger.LogDebug().Msg(fmt.Sprintf("Received message: %+v", data))
		fmt.Println("recieved message", data)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		var channel *db.ChannelModel
		if data.ChannelID != "" {
			channel, err = ws.srv.GetChannelById(data.ChannelID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			data.Recievers = channel.Participants()
		}
		data.SenderID = claims.ID
		messages <- data
	}
}
