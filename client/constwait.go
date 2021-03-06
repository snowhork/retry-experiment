package client

import (
	"context"
	"time"

	"github.com/snowhork/retry-experiment/server"
)

type ConstWait struct {
	WaitTime time.Duration
}

func (c *ConstWait) RequestWithRetry(ctx context.Context, s *server.Server, callback func(success bool)) error {
	for {
		success := s.Request(ctx)
		callback(success)

		if success {
			break
		} else {
			time.Sleep(c.WaitTime)
		}
	}

	return nil
}
