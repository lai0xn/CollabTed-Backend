package main

import (
	_ "github.com/CollabTed/CollabTed-Backend/docs"
	"github.com/CollabTed/CollabTed-Backend/internal/server"
	"github.com/CollabTed/CollabTed-Backend/pkg/redis"
	"github.com/CollabTed/CollabTed-Backend/prisma"
)

// @title			Squid Tech API
// @version		1.0
// @description	backend of the event management app.
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	s := server.NewServer(":8080")
	prisma.Connect()
	redis.Connect()
	s.Run()
}
