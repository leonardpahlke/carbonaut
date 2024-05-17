package logger_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"carbonaut.dev/pkg/util/logger"
)

func Example() {
	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr, logger.DefaultOptions)))

	slog.Info("Initializing")
	slog.Debug("Init done", "duration", 500*time.Millisecond)
	slog.Warn("Slow request!", "method", "GET", "path", "/api/users", "duration", 750*time.Millisecond)
	slog.Error("DB connection lost!", "err", errors.New("connection reset"), "db", "horalky")
	// Output:
}

func BenchmarkLog(b *testing.B) {
	b.StopTimer()
	l := slog.New(logger.NewHandler(os.Stderr, logger.DefaultOptions))

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		l.Info("benchmarking", "i", i)
		b.StopTimer()
	}
}
