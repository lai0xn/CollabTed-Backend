package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	ADDR string
}

type Server struct {
	config Config
}

func NewWithConfig(config Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Run() {
	r := chi.NewRouter()
	server := http.Server{
		Addr:              s.config.ADDR,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
