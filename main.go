package main

import (
	"context"
	"github.com/snowhork/retry-experiment/client"
	"github.com/snowhork/retry-experiment/server"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

type ClientType string

const (
	nowait             ClientType = "nowait"
	constTime          ClientType = "constTime"
	linearTimeout      ClientType = "linearTimeout"
	exponentialBackoff ClientType = "exponentialBackoff"
)

func do(clientType ClientType) error {
	ctx := context.Background()
	s := server.NewServer()
	eg, ctx := errgroup.WithContext(ctx)

	totalWorker := 100
	requestCountPerWorker := 200
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
		return err
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
	println("clientType: ", clientType)
	println("attempts:   ", allLog.TotalAttempts())
	println("success:    ", allLog.SuccessCount())
	println("duration:   ", delta.Milliseconds())

	return nil
}

func main() {
	for _, clientType := range []ClientType{
		//nowait,
		constTime,
		linearTimeout,
		exponentialBackoff,
	} {
		if err := do(clientType); err != nil {
			log.Fatal(err)
		}
	}
}
