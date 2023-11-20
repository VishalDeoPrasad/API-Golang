package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
		Handler:      API(),
	}
	log.Info().Str("port", api.Addr).Msg("main: API listening")

	api.ListenAndServe()
	api.Close()
	return nil
}

func Signup(c *gin.Context) {
	// Dummy implementation for the Signup handler
	c.JSON(200, gin.H{"message": "Signup Successful"})
}

func Login(c *gin.Context) {
	// Dummy implementation for the Login handler
	c.JSON(200, gin.H{"message": "Login Successful"})
}

func API() *gin.Engine {
	// Create a new Gin engine
	r := gin.New()

	// Use the middleware and recovery globally
	r.Use(gin.Logger(), gin.Recovery())

	// Define routes
	r.POST("/signup", Signup)
	r.POST("/login", Login)

	return r
}
