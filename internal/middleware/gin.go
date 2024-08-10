package middleware

import (
	"github.com/hedon954/go-matcher/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func WithRequestAndTrace() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set(response.RequestIDKey, c.GetHeader(response.XRequestID))
		c.Set(response.TraceIDKey, uuidV7())
	}
}

func uuidV7() string {
	return uuid.Must(uuid.NewV7()).String()
}
