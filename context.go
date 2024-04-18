package carthage

import (
	"context"
	"os"

	log "github.com/rs/zerolog"
)

type contextKey int

const (
	loggerContextKey = contextKey(iota)
)

func NewContextWithLogger(ctx context.Context, level log.Level) context.Context {
	log.SetGlobalLevel(level)
	logger := log.New(os.Stdout).With().Timestamp().Logger()
	return logger.WithContext(ctx)
}

func LoggerFromContext(ctx context.Context) *log.Logger {
	return log.Ctx(ctx)
}
