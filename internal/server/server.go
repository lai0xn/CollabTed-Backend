package server

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	ADDR string
	Log  *log.Logger
}

type Server struct {
	Addr   string
	logger *log.Logger
}

func NewWithConfig(config Config) *Server {
	return &Server{
		Addr:   config.ADDR,
		logger: config.Log,
	}
}

func (s *Server) Run() {
	r := chi.NewRouter()
	server := http.Server{
		Addr:              s.Addr,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	s.logger.Info("Server Started Listening", "port", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		s.logger.Error(err)
	}
}
