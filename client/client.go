package client

import (
	"context"
	"github.com/snowhork/retry-experiment/server"
	"time"
)

type Client interface {
	RequestWithRetry(ctx context.Context, s *server.Server, callback func(success bool)) error
}

type Worker struct {
	Client Client
	ID     int
	Need   int
	Log    Log
}

type Payload struct {
	WorkerID      int
	RequestNumber int
	AttemptCount  int
}

func (w *Worker) Request(ctx context.Context, s *server.Server) error {
	attemptCount := 0
	for i := 0; i < w.Need; i++ {
		if err := w.Client.RequestWithRetry(ctx, s, func(suc bool) {
			attemptCount += 1
			w.Log.rows = append(w.Log.rows,
				LogRow{
					Time:    time.Now(),
					Success: suc,
					Payload: Payload{
						WorkerID:      w.ID,
						RequestNumber: i + 1,
						AttemptCount:  attemptCount,
					},
				})
		}); err != nil {
			return err
		}
	}
	return nil
}
