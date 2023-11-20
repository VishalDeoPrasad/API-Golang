package handlers

import (
	"encoding/json"
	"golang/auth"
	"golang/middleware"
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
	a  *auth.Auth
}

func (h *handler) Signup(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TrackerIdKey).(string)
	if !ok{
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Track Id missing"})
		return
	}

	var nu models.NewUserReq
	body := c.Request.Body
	err := json.NewDecoder(body).Decode(&nu)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("Problem in reading request body for new user")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON format for new user"})
	}

	validate := validator.New()
	err = validate.Struct(nu)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("please provide name, email and password.")
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"msg": "Validation failed"})
	}

	err = h.db.AutoMigrate()
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("check handler")
	}

	newRecord, err := h.db.CreateUser(nu)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("check handler")
	}

	log.Info().Str("user_name", newRecord.Name).Str("Trace Id", traceId).Msg("User created successfully")
	c.JSON(http.StatusOK, newRecord)

}

func (h *handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TrackerIdKey).(string)
	if !ok{
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Track Id missing"})
	}
	var user models.LoginReq
	credit := c.Request.Body
	err := json.NewDecoder(credit).Decode(&user)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("Problem in reading request body for login")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON format for for login"})
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		log.Error().Err(err).Msg("please provide all credentials, Name and Password")
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"msg": "Validation failed, please provide all credentials"})
	}

	claims, err := h.db.Authentication(user)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("password didn't match: handler layer")
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"msg": "Authentication failed, check user and password or registered claims,"})
		return
	}

	// Generate a new token and put it in the Token field of the token struct
	// Define a new struct for the token
	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = h.a.GenerateToken(claims)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("generating token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "error in token generation"})
		return
	}

	log.Info().Str("user_name", user.Email).Msg("Login successfully...")
	c.JSON(http.StatusOK, gin.H{"Msg": "Login Succesfull..."})

}

func API(db *gorm.DB, a *auth.Auth) *gin.Engine {
	db_conn, err := service.NewConn(db)
	if err != nil {
		log.Error().Err(err).Msg("check handler")
	}

	h := handler{
		db: db_conn,
		a:  a,
	}

	// Create a new Gin engine
	r := gin.New()

	// Use the middleware and recovery globally
	r.Use(middleware.Logger(), gin.Recovery())

	// Define routes
	r.POST("/signup", h.Signup)

	r.POST("/login", h.Login)

	return r
}
