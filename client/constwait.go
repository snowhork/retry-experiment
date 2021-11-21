package client

import (
	"context"
	"time"

	"github.com/snowhork/retry-experiment/server"
)

type ConstWait struct {
	WaitTime time.Duration
	Log Log
}

func (c *ConstWait) RequestWithRetry(ctx context.Context, s *server.Server, payload int) error {
	success := false
	for !success {
		success = s.Request(ctx, payload)
		c.Log.rows = append(c.Log.rows, LogRow{
			Time: time.Now(),
			Payload: payload,
			Success: success,
		})

		if !success {
			time.Sleep(c.WaitTime)
		}
	}

	return nil
}

func (c *ConstWait) GetLog() *Log {
	return &c.Log
}
