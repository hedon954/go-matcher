package middleware

import (
	"github.com/hedon954/go-matcher/pkg/rand"
	"github.com/hedon954/go-matcher/pkg/response"

	"github.com/gin-gonic/gin"
)

func WithRequestAndTrace() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set(response.RequestIDKey, c.GetHeader(response.XRequestID))
		c.Set(response.TraceIDKey, rand.UUIDV7())
	}
}
