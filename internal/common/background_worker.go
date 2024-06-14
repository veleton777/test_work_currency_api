package common

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

func BackgroundWorker(ctx context.Context, delay time.Duration, logger *zerolog.Logger, fn func(ctx context.Context) error) {
	t := time.NewTicker(delay)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := fn(ctx); err != nil {
				logger.Err(err).Msg("background worker task run")
			}
		}
	}
}
