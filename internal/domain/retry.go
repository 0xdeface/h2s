package domain

import (
	"context"
	"log"
	"time"
)

func Retry(effector func(ctx context.Context) error, retries int, delay time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		for r := 0; ; r++ {
			err := effector(ctx)
			if err == nil || r >= retries {
				return err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
