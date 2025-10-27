package until

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

func Until(ctx context.Context, interval time.Duration, fn func() (bool, error)) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Warn().Msg("Operation timed out.")
			return ctx.Err()
		case <-ticker.C:
			success, err := fn()
			if err != nil {
				return err
			}

			if success {
				log.Debug().Msg("Function completed successfully.")
				return nil
			}

			log.Debug().Msg("Function not yet completed.")
		}
	}
}
