package main

import (
	"fmt"

	"github.com/CollabTed/CollabTed-Backend/internal/server"
	"github.com/CollabTed/CollabTed-Backend/pkg/logger"
	"github.com/spf13/viper"
)

func main() {
	// Read config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	s := server.NewWithConfig(server.Config{
		ADDR: viper.GetString("server.port"),
		Log:  logger.New("🌐 🗄️", false),
	})

	s.Run()
}
