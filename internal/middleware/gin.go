package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/hedon954/go-matcher/pkg/rand"
	"github.com/hedon954/go-matcher/pkg/response"
)

// WithRequestAndTrace is a middleware that sets the request id and trace id to the context.
// If the context has a span, it will use the trace id of the span as the trace id.
// Otherwise, it will generate a random trace id.
func WithRequestAndTrace() func(c *gin.Context) {
	return func(c *gin.Context) {
		var traceID string
		span := trace.SpanFromContext(c.Request.Context())
		if span != nil {
			traceID = span.SpanContext().TraceID().String()
		} else {
			traceID = rand.UUIDV7()
		}

		c.Set(response.RequestIDKey, c.GetHeader(response.XRequestID))
		c.Set(response.TraceIDKey, traceID)
	}
}
