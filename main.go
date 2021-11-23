package main

import (
	"context"
	"github.com/snowhork/retry-experiment/client"
	"github.com/snowhork/retry-experiment/server"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

var clientType = "nowait"
//var clientType = "const"

func do() error {
	ctx := context.Background()
	s := server.NewServer()
	eg, ctx := errgroup.WithContext(ctx)

	totalRequestCount := 1000
	clients := make([]client.Client, totalRequestCount)
	for i := 0; i < totalRequestCount; i++ {
		switch clientType {
		case "nowait":
			clients[i] = &client.NoWait{}
			break
		case "const":
			clients[i] = &client.ConstWait{WaitTime: 10*time.Millisecond}
			break
		}
	}

	for i := 0; i < totalRequestCount; i++ {
		i := i
		eg.Go(func() error {
			return clients[i].RequestWithRetry(ctx, s, i)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	var (
		allLog *client.Log
	)
	{
		logs := make([]*client.Log, totalRequestCount)
		for i, c := range clients {
			logs[i] = c.GetLog()
		}
		allLog = client.AggregateLog(logs)
	}

	delta := allLog.LatestTime().Sub(allLog.OldestTime())

	println("clientType: ", clientType)
	println("attempts:   ", allLog.TotalAttempts())
	println("success:    ", allLog.SuccessCount())
	println("duration:   ", delta.Milliseconds())


	return nil
}

func main() {
	if err := do(); err != nil {
		log.Fatal(err)
	}
}
