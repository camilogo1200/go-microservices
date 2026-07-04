package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func getMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	return mux
}

func getCors() *cors.Cors {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://*"},
		AllowedMethods: []string{
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions,
		},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		MaxAge:           700,
		//Debug: true
	})

	return c
}

func health(w http.ResponseWriter, r *http.Request) {
	now := fmt.Sprintf("Ok - %v" + time.Now().String())
	_, err := w.Write([]byte(now))
	if err != nil {
		return
	}
}
