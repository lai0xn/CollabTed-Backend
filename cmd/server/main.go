package main

import "github.com/CollabTed/CollabTed-Backend/internal/server"

func main() {
	s := server.NewWithConfig(server.Config{
		ADDR: ":8080",
	})

	s.Run()
}
