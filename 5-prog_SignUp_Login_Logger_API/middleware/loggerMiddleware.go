// * Zerolog logger middleware
package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type key string

const TrackerIdKey key = "1"

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		trackerId := uuid.NewString()
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, TrackerIdKey, trackerId)
		req := c.Request.WithContext(ctx)
		c.Request = req

		log.Info().Str("track ID", trackerId).Str("Method", c.Request.Method).
			Str("URL Path", c.Request.URL.Path).Msg("request started")

		defer log.Info().Str("track ID", trackerId).Str("Method", c.Request.Method).
			Str("URL Path", c.Request.URL.Path).Int("status Code", c.Writer.Status()).Msg("Request processing completed")

		c.Next()

	}
}
