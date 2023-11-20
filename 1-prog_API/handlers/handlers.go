package handlers

import (
	"github.com/gin-gonic/gin"
)

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
