package server

import (
	"log"
	"net/http"

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

	if err := http.ListenAndServe(s.config.ADDR, r); err != nil {
		log.Println(err)
	}
}
