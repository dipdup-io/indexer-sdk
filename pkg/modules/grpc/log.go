package grpc

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog/log"
)

func logCalls() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := log.With().Ctx(ctx).Str("module", "grpc").Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Fields(fields).Msg(msg)
		case logging.LevelInfo:
			l.Info().Fields(fields).Msg(msg)
		case logging.LevelWarn:
			l.Warn().Fields(fields).Msg(msg)
		case logging.LevelError:
			l.Error().Fields(fields).Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
