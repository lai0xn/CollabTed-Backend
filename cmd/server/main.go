package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		response := make(map[string]interface{})
		response["message"] = "Server is running"
		json.NewEncoder(w).Encode(response)
	})
	log.Println("Server is running")
	http.ListenAndServe(":8080", r)
}
