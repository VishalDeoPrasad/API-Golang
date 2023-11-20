package middleware

import (
	"context"
	"errors"
	"golang/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (m *Mid) Authenticate(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceId, ok := ctx.Value(TrackerIdKey).(string)
		if !ok {
			log.Error().Msg("trace id not present in the context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			return
		}

		// Getting the Authorization header
		authHeader := c.Request.Header.Get("Authorization")

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			err := errors.New("expected authorization header format: Bearer <token>")
			log.Error().Err(err).Str("Trace Id", traceId).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// ValidateToken presumably checks the token for validity and returns claims if it's valid
		claims, err := m.a.ValidateToken(parts[1])
		// If there is an error, log it and return an Unauthorized error message
		if err != nil {
			log.Error().Err(err).Str("Trace Id", traceId).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
			return
		}

		// If the token is valid, then add it to the context
		ctx = context.WithValue(ctx, auth.Key, claims)

		// Creates a new request with the updated context and assign it back to the gin context
		req := c.Request.WithContext(ctx)
		c.Request = req

		// Proceed to the next middleware or handler function
		next(c)
	}
}
