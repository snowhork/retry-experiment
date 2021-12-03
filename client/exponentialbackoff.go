package client

import (
	"context"
	"time"

	"github.com/snowhork/retry-experiment/server"
)

type ExponentialBackoff struct {
	WaitTime time.Duration
}

func (c *ExponentialBackoff) RequestWithRetry(ctx context.Context, s *server.Server, callback func(success bool)) error {
	waitTime := c.WaitTime
	mult := 1

	for {
		success := s.Request(ctx)
		callback(success)

		if success {
			break
		} else {
			time.Sleep(time.Duration(mult) * waitTime)
		}
		mult *= 2
	}

	return nil
}
