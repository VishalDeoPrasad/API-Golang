package main

import (
	"context"
	"fmt"
	"golang/auth"
	"golang/database"
	"golang/handlers"
	"golang/middleware"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

func main() {
	err := startApp()
	if err != nil {
		log.Panic().Err(err).Send()
	}
}

func startApp() error {
	//=========================================================================
	//read keys for authentication
	privatePEM, err := os.ReadFile("private.pem")
	//This is an error check. If there was an error in reading the file
	if err != nil {
		return fmt.Errorf("reading auth private key %w", err)
		//it returns an error message using fmt.
		//The %w verb is used to wrap the original error within a new error message for context.
	}
	//PEM is a text-based encoding format that's commonly used for various types of data, including cryptographic keys.
	//converting it into a format that can be used for cryptographic operations, such as signing or decrypting data.

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("parsing auth private key %w", err)
	}

	publicPEM, err := os.ReadFile("pubkey.pem")
	if err != nil {
		return fmt.Errorf("reading auth public key %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return fmt.Errorf("parsing auth public key %w", err)
	}
	//received the instance of Auth struct and a nil error
	a, err := auth.NewAuth(privateKey, publicKey)
	if err != nil {
		return fmt.Errorf("constructing auth %w", err)
	}

	m, err := middleware.NewMid(a)
	if err != nil {
		return fmt.Errorf("constructing mid %w", err)
	}

	// =========================================================================
	// Start Database
	log.Info().Msg("main : Started : Initializing db support")
	db, err := database.Open()
	if err != nil {
		return fmt.Errorf("connecting to db %w", err)
	}
	pg, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = pg.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("database is not connected: %w ", err)
	}
	//===========================================================================
	// Initialize http service
	api := http.Server{
		Addr:         ":8080",
		ReadTimeout:  8000 * time.Second,
		WriteTimeout: 800 * time.Second,
		IdleTimeout:  800 * time.Second,
		Handler:      handlers.API(db,a, m),
	}
	log.Info().Str("port", api.Addr).Msg("main: API listening")

	api.ListenAndServe()
	api.Close()
	return nil
}
