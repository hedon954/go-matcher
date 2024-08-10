package main

import (
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const name = "github.com/hedon954/go-matcher/example/otel"

var (
	tracer  = otel.Tracer(name)
	meter   = otel.Meter(name)
	logger  = otelslog.NewLogger(name)
	rollCnt metric.Int64Counter
)

func init() {
	var err error
	rollCnt, err = meter.Int64Counter("dice.rolls",
		metric.WithDescription("The number of rolls by roll value"),
		metric.WithUnit("{roll}"),
	)
	if err != nil {
		panic(err)
	}

}

func rollDice(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "roll")
	defer span.End()

	roll := 1 + rand.IntN(6)

	var msg string
	if player := r.PathValue("player"); player != "" {
		msg = fmt.Sprintf("%s is rolling the dice", player)
	} else {
		msg = "An anonymous player is rolling the dice"
	}
	logger.InfoContext(ctx, msg, "result", roll)

	rollValueAttr := attribute.Int("roll.value", roll)
	rollMsgAttr := attribute.String("roll.message", msg)
	span.SetAttributes(rollValueAttr, rollMsgAttr)
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr, rollMsgAttr))

	if rand.IntN(10) > 8 {
		handleRollDice(ctx)
	}

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func handleRollDice(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "roll.handle")
	defer span.End()
	rollMsgAttr := attribute.String("msg", "this is a message")
	span.SetAttributes(rollMsgAttr)
}
