package postgres

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

type logQueryHook struct{}

// BeforeQuery -
func (h *logQueryHook) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	event.StartTime = time.Now()
	return ctx, nil
}

func (h *logQueryHook) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	if err != nil {
		return err
	}

	// log.Trace().Interface("params", event.Params).Msg("")
	if event.Err != nil {
		log.Trace().Msgf("[%d ms] %s : %s", time.Since(event.StartTime).Milliseconds(), event.Err.Error(), string(query))
	} else {
		log.Trace().Msgf("[%d ms] %d rows | %s", time.Since(event.StartTime).Milliseconds(), event.Result.RowsReturned(), string(query))
	}
	return nil
}
