package client

import (
	"context"
	"time"

	"github.com/snowhork/retry-experiment/server"
)

type NoWait struct {
	Log Log
}

func (c *NoWait) RequestWithRetry(ctx context.Context, s *server.Server, payload int) error {
	success := false
	for !success {
		success = s.Request(ctx, payload)
		c.Log.rows = append(c.Log.rows, LogRow{
			Time: time.Now(),
			Payload: payload,
			Success: success,
		})

	}

	return nil
}

func (c *NoWait) GetLog() *Log {
	return &c.Log
}
