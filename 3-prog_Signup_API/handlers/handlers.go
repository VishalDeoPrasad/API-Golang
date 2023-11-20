package handlers

import (
	"encoding/json"
	"golang/models"
	"golang/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type handler struct {
	db *service.Conn
}

func (h *handler) Signup(c *gin.Context) {
	var nu models.NewUserReq
	body := c.Request.Body
	err := json.NewDecoder(body).Decode(&nu)
	if err != nil {
		log.Error().Err(err).Msg("Problem in reading request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON format"})
	}

	validate := validator.New()
	err = validate.Struct(nu)
	if err != nil {
		log.Error().Err(err).Msg("please provide name, email and password.")
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"msg": "Validation failed"})
	}

	err = h.db.AutoMigrate()
	if err != nil {
		log.Error().Err(err).Msg("check handler")
	}

	newRecord, err := h.db.CreateUser(nu)
	if err != nil {
		log.Error().Err(err).Msg("check handler")
	}

	log.Info().Str("user_name", newRecord.Name).Msg("User created successfully")
	c.JSON(http.StatusOK, newRecord)

}

func (h *handler) Login(c *gin.Context) {
	// Dummy implementation for the Login handler
	c.JSON(200, gin.H{"message": "Login Successful"})
}

func API(db *gorm.DB) *gin.Engine {
	db_conn, err := service.NewConn(db)
	if err != nil {
		log.Error().Err(err).Msg("check handler")
	}

	h := handler{
		db: db_conn,
	}

	// Create a new Gin engine
	r := gin.New()

	// Use the middleware and recovery globally
	r.Use(gin.Logger(), gin.Recovery())

	// Define routes
	r.POST("/signup", h.Signup)

	r.POST("/login", h.Login)

	return r
}
