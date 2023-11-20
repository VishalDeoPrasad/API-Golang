package main

import (
	"golang/handlers"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	err := startApp()
	if err != nil {
		log.Panic().Err(err).Send()
	}
}

func startApp() error {
	// Initialize http service
	api := http.Server{
		Addr:         ":8080",
		ReadTimeout:  8000 * time.Second,
		WriteTimeout: 800 * time.Second,
		IdleTimeout:  800 * time.Second,
		Handler:      handlers.API(),
	}
	log.Info().Str("port", api.Addr).Msg("main: API listening")
	
	api.ListenAndServe()
	api.Close()
	return nil
}
