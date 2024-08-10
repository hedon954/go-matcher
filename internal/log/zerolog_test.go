package log

import (
	"context"
	"testing"

	"github.com/hedon954/go-matcher/pkg/response"
	"github.com/stretchr/testify/assert"
)

// nolint:staticcheck
func TestLog(t *testing.T) {
	ctx := context.WithValue(context.Background(), response.RequestIDKey, "request_id")
	ctx = context.WithValue(ctx, response.TraceIDKey, "trace_id")
	Ctx(ctx).Info().Msg("with ctx")
	Info().Msg("info")
	Debug().Msg("debug")
	Error().Msg("error")
	Trace().Msg("trace")
	Warn().Msg("warn")
	assert.Panics(t, func() {
		Panic().Msg("panic")
	})
}
