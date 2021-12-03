package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/snowhork/retry-experiment/client"
	"github.com/snowhork/retry-experiment/server"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"sort"
	"time"
)

type ClientType string

const (
	nowait                            ClientType = "nowait"
	constTime                         ClientType = "constTime"
	linearTimeout                     ClientType = "linearTimeout"
	exponentialBackoff                ClientType = "exponentialBackoff"
	exponentialBackoffWithConstJitter ClientType = "exponentialBackoffWithConstJitter"
	exponentialBackoffWithHalfJitter  ClientType = "exponentialBackoffWithHalfJitter"
	exponentialBackoffWithFullJitter  ClientType = "exponentialBackoffWithFullJitter"
)

type summary struct {
	clientType   ClientType
	totalWorker  int
	attempts     int
	successCount int
	successRate  float32
	tryRate      float32
	duration     int64
}

func do(now int64, clientType ClientType, totalWorker int) (summary, error) {
	ctx := context.Background()
	s := server.NewServer()
	eg, ctx := errgroup.WithContext(ctx)

	requestCountPerWorker := 10
	workers := make([]client.Worker, totalWorker)
	for i := 0; i < totalWorker; i++ {
		var c client.Client
		switch clientType {
		case nowait:
			c = &client.NoWait{}
		case constTime:
			c = &client.ConstWait{WaitTime: 1 * time.Millisecond}
			break
		case linearTimeout:
			c = &client.LinearTimeout{WaitTime: 1 * time.Millisecond}
			break
		case exponentialBackoff:
			c = &client.ExponentialBackoff{WaitTime: 1 * time.Millisecond}
			break
		case exponentialBackoffWithConstJitter:
			c = &client.ExponentialBackoffWithConstJitter{WaitTime: 1 * time.Millisecond}
			break
		case exponentialBackoffWithHalfJitter:
			c = &client.ExponentialBackoffWithHalfJitter{WaitTime: 1 * time.Millisecond}
			break
		case exponentialBackoffWithFullJitter:
			c = &client.ExponentialBackoffWithFullJitter{WaitTime: 1 * time.Millisecond}
			break
		}

		workers[i] = client.Worker{
			Client: c,
			ID:     i + 1,
			Need:   requestCountPerWorker,
		}
	}

	for i := 0; i < totalWorker; i++ {
		i := i
		eg.Go(func() error {
			return workers[i].Request(ctx, s)
		})
	}

	if err := eg.Wait(); err != nil {
		return summary{}, err
	}

	var (
		allLog *client.Log
	)
	{
		logs := make([]client.Log, totalWorker)
		for i, w := range workers {
			logs[i] = w.Log
		}
		allLog = client.AggregateLog(logs)
	}

	delta := allLog.LatestTime().Sub(allLog.OldestTime())

	println("----------")
	println("clientType:      ", clientType)
	println("attempts:        ", allLog.TotalAttempts())
	println("success:         ", allLog.SuccessCount())
	println("success rate(%): ", 100*allLog.SuccessCount()/allLog.TotalAttempts())
	println("try rate:        ", 1.0*allLog.TotalAttempts()/allLog.SuccessCount())
	println("duration:        ", delta.Milliseconds())

	f, err := os.Create(fmt.Sprintf("./%d_%s_%d.csv", now, clientType, totalWorker))
	if err != nil {
		return summary{}, err
	}
	defer f.Close()

	sort.Slice(s.Counter.Logs, func(i, j int) bool {
		return s.Counter.Logs[i].Time < s.Counter.Logs[j].Time
	})

	writer := csv.NewWriter(f)
	if err := writer.Write([]string{"time", "counter"}); err != nil {
		return summary{}, err
	}
	for _, l := range s.Counter.Logs {
		if err := writer.Write([]string{fmt.Sprintf("%d", l.Time), fmt.Sprintf("%d", l.Counter)}); err != nil {
			return summary{}, err
		}
	}

	writer.Flush()

	return summary{
		clientType:   clientType,
		totalWorker:  totalWorker,
		attempts:     allLog.TotalAttempts(),
		successCount: allLog.SuccessCount(),
		successRate:  100.0 * float32(allLog.SuccessCount()) / float32(allLog.TotalAttempts()),
		tryRate:      float32(allLog.TotalAttempts()) / float32(allLog.SuccessCount()),
		duration:     delta.Milliseconds(),
	}, nil
}

func workerLoop(now int64, clientType ClientType) error {
	f, err := os.Create(fmt.Sprintf("./summary_%d_%s.csv", now, clientType))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	if err := writer.Write([]string{"worker", "attempts", "successCount", "successRate", "tryRate", "duration"}); err != nil {
		return err
	}

	for _, worker := range []int{10, 15, 20, 25, 30, 35} {
		s, err := do(now, clientType, worker)
		if err != nil {
			return err
		}

		fmt.Printf("%v\n", s)

		if err := writer.Write([]string{
			fmt.Sprintf("%d", s.totalWorker),
			fmt.Sprintf("%d", s.attempts),
			fmt.Sprintf("%d", s.successCount),
			fmt.Sprintf("%f", s.successRate),
			fmt.Sprintf("%f", s.tryRate),
			fmt.Sprintf("%d", s.duration),
		}); err != nil {
			return err
		}
	}

	writer.Flush()
	println(f.Name())

	return nil
}

func main() {
	now := time.Now().UnixNano()
	println(now)
	for _, clientType := range []ClientType{
		//nowait,
		constTime,
		//linearTimeout,
		exponentialBackoff,
		exponentialBackoffWithConstJitter,
		exponentialBackoffWithHalfJitter,
		exponentialBackoffWithFullJitter,
	} {
		if err := workerLoop(now, clientType); err != nil {
			log.Fatal(err)
		}
	}
}
