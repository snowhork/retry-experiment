package server

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

type Server struct {
	semaphore *semaphore.Weighted
	Counter   Counter
}

type Counter struct {
	mu    sync.Mutex
	value int
	Logs  []CounterLog
}

type CounterLog struct {
	Counter int
	Time    int64
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	c.Logs = append(c.Logs, CounterLog{Counter: c.value, Time: time.Now().UnixNano()})
}

func (c *Counter) Decrement() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value--
	c.Logs = append(c.Logs, CounterLog{Counter: c.value, Time: time.Now().UnixNano()})
}

func NewServer() *Server {
	return &Server{
		semaphore: semaphore.NewWeighted(3),
	}
}

func (s *Server) Request(ctx context.Context) bool {
	s.Counter.Increment()
	defer s.Counter.Decrement()

	if s.semaphore.TryAcquire(1) {
		defer func() {
			s.semaphore.Release(1)
		}()

		//millsec := 10 + rand.Intn(10)

		millsec := 3
		time.Sleep(time.Duration(millsec) * time.Millisecond)

		return true
	}
	//millsec := 5 + rand.Intn(5)
	//millsec := 5
	//Time.Sleep(Time.Duration(millsec)*Time.Millisecond)

	//Time.Sleep(5*Time.Millisecond)
	return false
}
