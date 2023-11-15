package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"service-app/auth"
	"service-app/middlewares"
	"service-app/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type handler struct {
	s models.Store
	a *auth.Auth
}

// Signup is a method for the handler struct which handles user registration
func (h *handler) Signup(c *gin.Context) {
	//c is an object that represents the current HTTP request and response.
	/*
		This line is extracting the context (ctx) from the HTTP request (c.Request) using the Context method.
		Let's break down what this means:
	*/
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	//the context is later used to retrieve the traceId.
	if !ok {
		// If the traceId isn't found in the request, log an error and return
		log.Error().Msg("traceId missing from context")

		// You can create a custom error message or use a predefined one
		err := errors.New("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Define a NewUser variable
	var nu models.NewUser

	// Attempt to decode JSON from the request body into the NewUser variable
	err := json.NewDecoder(c.Request.Body).Decode(&nu)
	if err != nil {
		// If there is an error in decoding, log the error and return
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON format"})
		return
	}

	// Create a new validator and validate the NewUser variable
	/*
		The validator package uses struct tags (tags in the struct definition) to define validation 
		rules for struct fields. For example:*/
	validate := validator.New()
	err = validate.Struct(nu)
	if err != nil {
		// If validation fails, log the error and return
		
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Validation failed"})
		return
	}

	// Attempt to create the user
	usr, err := h.s.CreateUser(ctx, nu)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("user signup problem")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "user signup failed"})
		return
	}

	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, usr)
}

// Login is a method for the handler struct which handles user login
func (h *handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// Define a new struct for login data
	var login struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	// Attempt to decode JSON from the request body into the login variable
	err := json.NewDecoder(c.Request.Body).Decode(&login)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// Create a new validator and validate the login variable
	validate := validator.New()
	err = validate.Struct(login)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide Email and Password"})
		return
	}

	// Attempt to authenticate the user with the email and password
	claims, err := h.s.Authenticate(ctx, login.Email, login.Password)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "login failed"})
		return
	}

	// Define a new struct for the token
	var tkn struct {
		Token string `json:"token"`
	}

	// Generate a new token and put it in the Token field of the token struct
	tkn.Token, err = h.a.GenerateToken(claims)
	if err != nil {
		log.Error().Err(err).Msg("generating token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// If everything goes right, respond with the token
	c.JSON(http.StatusOK, tkn)

}

func (h *handler) AddInventory(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	claims, ok := ctx.Value(auth.Key).(jwt.RegisteredClaims)
	if !ok {
		log.Error().Str("Trace Id", traceId).Msg("login first")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	var newInv models.NewInventory
	err := json.NewDecoder(c.Request.Body).Decode(&newInv)
	if err != nil {
		log.Error().Str("Trace Id", traceId).Msg("login first")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}
	validate := validator.New()
	err = validate.Struct(newInv)

	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"msg": "please provide Item Name and Quantity"})
		return
	}
	uid, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.
			StatusInternalServerError)})
		return
	}
	inv, err := h.s.CreatInventory(ctx, newInv, uint(uid))
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Inventory creation failed"})
		return
	}

	c.JSON(http.StatusOK, inv)

}

func (h *handler) ViewInventory(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	claims, ok := ctx.Value(auth.Key).(jwt.RegisteredClaims)
	if !ok {
		log.Error().Str("Trace Id", traceId).Msg("login first")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}
	shirts, total, err := h.s.ViewInventory(ctx, claims.Subject)

	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "problem in viewing inventory"})
		return
	}
	m := gin.H{"inv": shirts, "total_cost": total}
	c.JSON(http.StatusOK, m)
}
