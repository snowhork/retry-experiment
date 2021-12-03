package client

import (
	"context"
	"math/rand"
	"time"

	"github.com/snowhork/retry-experiment/server"
)

type ExponentialBackoffWithHalfJitter struct {
	WaitTime time.Duration
}

func (c *ExponentialBackoffWithHalfJitter) RequestWithRetry(ctx context.Context, s *server.Server, callback func(success bool)) error {
	waitTime := c.WaitTime
	mult := 1

	for {
		success := s.Request(ctx)
		callback(success)

		if success {
			break
		} else {
			base := time.Duration(mult) * waitTime
			jitter := time.Duration(rand.Float32() * float32(base/2))
			time.Sleep(base/2 + jitter)
		}
		mult *= 2
	}

	return nil
}
