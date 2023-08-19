package grpc

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func logCalls() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := log.With().Ctx(ctx).Str("module", "grpc").Logger()

		event := new(zerolog.Event)
		switch lvl {
		case logging.LevelDebug:
			event = l.Debug()
		case logging.LevelInfo:
			event = l.Info()
		case logging.LevelWarn:
			event = l.Warn()
		case logging.LevelError:
			event = l.Error()
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}

		event.Fields(fields).Msg(msg)
	})
}
