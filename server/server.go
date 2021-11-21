package server

import (
	"context"
	"golang.org/x/sync/semaphore"
	"time"
)

type Server struct {
	semaphore    *semaphore.Weighted
}

func NewServer() *Server {
	return &Server{
		semaphore: semaphore.NewWeighted(10),
	}
}

func (s *Server) Request(ctx context.Context, payload int) bool {
	if s.semaphore.TryAcquire(1) {
		defer s.semaphore.Release(1)
		time.Sleep(10*time.Millisecond)

		return true
	}

	return false
}
