package client

import (
	"context"
	"math/rand"
	"time"

	"github.com/snowhork/retry-experiment/server"
)

type ExponentialBackoffWithFullJitter struct {
	WaitTime time.Duration
}

func (c *ExponentialBackoffWithFullJitter) RequestWithRetry(ctx context.Context, s *server.Server, callback func(success bool)) error {
	waitTime := c.WaitTime
	mult := 1

	for {
		success := s.Request(ctx)
		callback(success)

		if success {
			break
		} else {
			jitter := time.Duration(rand.Float32() * float32(time.Duration(mult)*waitTime))
			time.Sleep(jitter)
		}
		mult *= 2
	}

	return nil
}
