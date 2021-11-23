package client

import (
	"context"
	"github.com/snowhork/retry-experiment/server"
)

type Client interface {
	RequestWithRetry(ctx context.Context, s *server.Server, payload int) error
	GetLog() *Log
}
