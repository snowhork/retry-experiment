package client

import (
	"context"
	"github.com/snowhork/retry-experiment/server"
	"time"
)

type Client interface {
	RequestWithRetry(ctx context.Context, s *server.Server, payload int) error
	GetLog() *Log
}

type Log struct {
	rows []LogRow
}

type LogRow struct {
	Time time.Time
	Payload int
	Success bool
}

func (l *Log) TotalAttempts() int {
	return len(l.rows)
}

func (l *Log) SuccessCount() int {
	res := 0
	for _, row := range l.rows {
		if row.Success {
			res += 1
		}
	}

	return res
}