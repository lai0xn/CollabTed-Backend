package main

import (
	"github.com/CollabTed/CollabTed-Backend/internal/server"
	"github.com/CollabTed/CollabTed-Backend/pkg/logger"
)

func main() {
	s := server.NewWithConfig(server.Config{
		ADDR: ":8080",
		Log:  logger.New("🌐 🗄️", false),
	})

	s.Run()
}
