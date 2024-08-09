package zerolog

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
)

/*
BenchmarkZeroLog_simple-12    	16250253	      72.47 ns/op
BenchmarkZeroLog_normal-12    	  634918	      1874 ns/op
BenchmarkZeroLog_large-12     	  214906	      5550 ns/op

BenchmarkSlog_simple-12       	 1514552	      795.7 ns/op
BenchmarkSlog_normal-12       	  318796	      3655 ns/op
BenchmarkSlog_large-12        	  116523	      10231 ns/op
*/

func BenchmarkZeroLog_simple(b *testing.B) {
	log := zerolog.New(io.Discard)
	for i := 0; i < b.N; i++ {
		log.Error().Msg("this is a error message this is a error message this is a error message this is a error message")
	}
}

func BenchmarkZeroLog_normal(b *testing.B) {
	log := zerolog.New(io.Discard)
	for i := 0; i < b.N; i++ {
		log.Error().
			Int64("id", 1000).
			Float64("f64", 3.14).
			Str("str", "this is a string info").
			Any("map", map[string]interface{}{
				"key":   "value",
				"int":   123,
				"float": 3.14,
				"bool":  true,
				"slice": []string{"a", "b", "c"},
				"map":   map[string]string{"a": "b", "c": "d"},
				"err":   errors.New("error"),
			}).
			Msg("this is a error message")
	}
}

func BenchmarkZeroLog_large(b *testing.B) {
	log := zerolog.New(io.Discard)
	for i := 0; i < b.N; i++ {
		log.Error().
			Int64("id", 1000).
			Float64("f64", 3.14).
			Str("str", "this is a string info").
			Int64("id", 1000).
			Float64("f64", 3.14).
			Str("str", "this is a string info").
			Int64("id", 1000).
			Float64("f64", 3.14).
			Str("str", "this is a string info").
			Int64("id", 1000).
			Float64("f64", 3.14).
			Str("str", "this is a string info").
			Any("map", map[string]interface{}{
				"key":   "value",
				"int":   123,
				"float": 3.14,
				"bool":  true,
				"slice": []string{"a", "b", "c"},
				"map":   map[string]string{"a": "b", "c": "d"},
				"err":   errors.New("error"),
			}).
			Any("map", map[string]interface{}{
				"key":   "value",
				"int":   123,
				"float": 3.14,
				"bool":  true,
				"slice": []string{"a", "b", "c"},
				"map":   map[string]string{"a": "b", "c": "d"},
				"err":   errors.New("error"),
			}).
			Any("map", map[string]interface{}{
				"key":   "value",
				"int":   123,
				"float": 3.14,
				"bool":  true,
				"slice": []string{"a", "b", "c"},
				"map":   map[string]string{"a": "b", "c": "d"},
				"err":   errors.New("error"),
			}).
			Msg("this is a error message")
	}
}

func BenchmarkSlog_simple(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	for i := 0; i < b.N; i++ {
		logger.Info("this is a string info this is a error message this is a error message this is a error message")
	}
}

func BenchmarkSlog_normal(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	for i := 0; i < b.N; i++ {
		logger.Info("this is a string info",
			slog.Int64("id", 1000),
			slog.Float64("f64", 3.14),
			slog.String("str", "this is a string info"),
			slog.Any("map", map[string]interface{}{
				"key":   "value",
				"int":   123,
				"float": 3.14,
				"bool":  true,
				"slice": []string{"a", "b", "c"},
				"map":   map[string]string{"a": "b", "c": "d"},
				"err":   errors.New("error"),
			}),
		)
	}
}

func BenchmarkSlog_large(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	for i := 0; i < b.N; i++ {
		logger.Info("this is a string info",
			slog.Int64("id", 1000),
			slog.Float64("f64", 3.14),
			slog.String("str", "this is a string info"),
			slog.Int64("id", 1000),
			slog.Float64("f64", 3.14),
			slog.String("str", "this is a string info"),
			slog.Int64("id", 1000),
			slog.Float64("f64", 3.14),
			slog.String("str", "this is a string info"),
			slog.Int64("id", 1000),
			slog.Float64("f64", 3.14),
			slog.String("str", "this is a string info"),
			slog.Any("map", map[string]interface{}{
				"key":   "value",
				"int":   123,
				"float": 3.14,
				"bool":  true,
				"slice": []string{"a", "b", "c"},
				"map":   map[string]string{"a": "b", "c": "d"},
				"err":   errors.New("error"),
			}),
			slog.Any("map", map[string]interface{}{
				"key":   "value",
				"int":   123,
				"float": 3.14,
				"bool":  true,
				"slice": []string{"a", "b", "c"},
				"map":   map[string]string{"a": "b", "c": "d"},
				"err":   errors.New("error"),
			}),
			slog.Any("map", map[string]interface{}{
				"key":   "value",
				"int":   123,
				"float": 3.14,
				"bool":  true,
				"slice": []string{"a", "b", "c"},
				"map":   map[string]string{"a": "b", "c": "d"},
				"err":   errors.New("error"),
			}),
		)
	}
}
