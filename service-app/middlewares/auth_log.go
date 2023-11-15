package middlewares

import (
	"context"
	"errors"
	"net/http"
	"service-app/auth"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type key string

const TraceIdKey key = "1"

// Mid is a structure that holds an authenticated session.
// This is typically used for maintaining user sessions or secure transactions.
type Mid struct {
	// 'a' attribute is a pointer to an 'Auth' object.
	// It's important to note that 'a'
	//is a pointer because we want to refer to the original 'Auth' object and not a COPY of it.
	a *auth.Auth
}

// NewMid is a function which takes an 'Auth' object pointer
// and returns a Mid instance and an error.
// Purpose of this function is to initialize
// and return a new instance of 'Mid' structure.
func NewMid(a *auth.Auth) (Mid, error) {
	// It first checks if 'a' is nil
	// 'a' should not be nil because 'nil' indicates that the 'Auth' object does not exist.
	if a == nil {
		// An error is returned when 'a' is 'nil'.
		return Mid{}, errors.New("auth can't be nil")
	}
	//If 'a' is not 'nil', a new 'Mid' instance is returned with 'a' as a field.
	// A nil error is returned, indicating that there were no issues with the initialization.
	return Mid{a: a}, nil
}

func (m *Mid) Log() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Generate a new unique identifier (UUID)
		traceId := uuid.NewString()

		// Fetch the current context from the gin context
		ctx := c.Request.Context()

		// Add the trace id in context so it can be used by upcoming processes in this request's lifecycle
		ctx = context.WithValue(ctx, TraceIdKey, traceId)

		// The 'WithContext' method on 'c.Request' creates a new copy of the request ('req'),
		// but with an updated context ('ctx') that contains our trace ID.
		// The original request does not get changed by this; we're simply creating a new version of it ('req').
		req := c.Request.WithContext(ctx)

		// Now, we want to carry forward this updated request (that has the new context) through our application.
		// So, we replace 'c.Request' (the original request) with 'req' (the new version with the updated context).
		// After this line, when we use 'c.Request' in this function or pass it to others, it'll be this new version
		// that carries our trace ID in its context.
		c.Request = req

		log.Info().Str("Trace Id", traceId).Str("Method", c.Request.Method).
			Str("URL Path", c.Request.URL.Path).Msg("request started")
		// After the request is processed by the next handler, logs the info again with status code
		defer log.Info().Str("Trace Id", traceId).Str("Method", c.Request.Method).
			Str("URL Path", c.Request.URL.Path).
			Int("status Code", c.Writer.Status()).Msg("Request processing completed")

		//we use c.Next only when we are using r.Use() method to assign middlewares
		c.Next()
	}
}

// Authenticate is a method that defines a Middleware function for gin HTTP framework
func (m *Mid) Authenticate(next gin.HandlerFunc) gin.HandlerFunc {
	// This middleware function is returned
	return func(c *gin.Context) {
		// We get the current request context
		ctx := c.Request.Context()

		// Extract the traceId from the request context
		// We assert the type to string since context.Value returns an interface{}
		traceId, ok := ctx.Value(TraceIdKey).(string)

		// If traceId not present then log the error and return an error message
		// ok is false if the type assertion was not successful
		if !ok {
			// Using a structured logging package (zerolog) to log the error
			log.Error().Msg("trace id not present in the context")

			// Sending error response using gin context
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			return
		}

		// Getting the Authorization header
		authHeader := c.Request.Header.Get("Authorization")

		// Splitting the Authorization header based on the space character.
		// Boats "Bearer" and the actual token
		parts := strings.Split(authHeader, " ")
		// Checking the format of the Authorization header
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			// If the header format doesn't match required format, log and send an error
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
