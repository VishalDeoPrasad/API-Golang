package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	r := gin.Default()

	r.POST("/signup", SignupHandler)

	// Run the server
	r.Run(":8080")
}

func SignupHandler(c *gin.Context) {
	// Extract context from the request
	ctx := c.Request.Context()

	// Dummy middleware to add traceId to the context (in real scenario, middleware might add actual traceId)
	ctx = AddTraceIdToContext(ctx, "123456789")

	// Retrieve traceId from the context
	traceId, ok := ctx.Value("traceId").(string)
	if !ok {
		// If the traceId isn't found in the request, log an error and return
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	fmt.Println("TraceId:", traceId)

	// Rest of the SignupHandler logic...
}

func AddTraceIdToContext(ctx context.Context, traceId string) context.Context {
	// Dummy middleware to add traceId to the context
	return context.WithValue(ctx, "traceId", traceId)
}
