package server

import (
	"context"
	"golang.org/x/sync/semaphore"
	"time"
)

type Server struct {
	semaphore *semaphore.Weighted
}

func NewServer() *Server {
	return &Server{
		semaphore: semaphore.NewWeighted(10),
	}
}

func (s *Server) Request(ctx context.Context) bool {
	if s.semaphore.TryAcquire(1) {
		defer s.semaphore.Release(1)

		//millsec := 10 + rand.Intn(10)
		millsec := 1
		time.Sleep(time.Duration(millsec) * time.Millisecond)

		return true
	}
	//millsec := 5 + rand.Intn(5)
	//millsec := 5
	//time.Sleep(time.Duration(millsec)*time.Millisecond)

	//time.Sleep(5*time.Millisecond)
	return false
}
