package main

import (
	_ "github.com/CollabTED/CollabTed-Backend/docs"
	"github.com/CollabTED/CollabTed-Backend/internal/server"
	"github.com/CollabTED/CollabTed-Backend/internal/ws"
	"github.com/CollabTED/CollabTed-Backend/pkg/cloudinary"
	"github.com/CollabTED/CollabTed-Backend/pkg/redis"
	"github.com/CollabTED/CollabTed-Backend/prisma"
)

// @title			CollabTED
// @version		1.0
// @description	REST Api of the CollabTED project.
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	redis.Connect()
	cloudinary.Connect()
	s := server.NewServer(":8080")
	prisma.Connect()
	go ws.Hub()
	s.Run()
}
