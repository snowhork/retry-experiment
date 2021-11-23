package client

import (
	"context"
	"github.com/snowhork/retry-experiment/server"
)

type NoWait struct {
}

func (c *NoWait) RequestWithRetry(ctx context.Context, s *server.Server, callback func(success bool)) error {
	for {
		success := s.Request(ctx)
		callback(success)

		if success {
			break
		}
	}

	return nil
}
