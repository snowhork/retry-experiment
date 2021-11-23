package client

import (
	"context"
	"time"

	"github.com/snowhork/retry-experiment/server"
)

type LinearTimeout struct {
	WaitTime time.Duration
}

func (c *LinearTimeout) RequestWithRetry(ctx context.Context, s *server.Server, callback func(success bool)) error {
	attempts := 1

	for {
		success := s.Request(ctx)
		callback(success)

		if success {
			break
		} else {
			time.Sleep(time.Duration(attempts) * c.WaitTime)
		}
	}

	return nil
}
