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
			clients[i] = &client.ConstWait{WaitTime: time.Millisecond}
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

	//if totalRequestCount != s.SuccessCount {
	//	return fmt.Errorf(
	//		"total requst count expected to be %d, but got %d",
	//		totalRequestCount,
	//		s.SuccessCount)
	//}

	totalAttempts := 0
	for _, c := range clients {
		totalAttempts += c.GetLog().TotalAttempts()
	}

	println("attempts: ", totalAttempts)

	totalSuccess := 0
	for _, c := range clients {
		totalSuccess += c.GetLog().SuccessCount()
	}

	println("success: ", totalSuccess)

	return nil
}

func main() {
	if err := do(); err != nil {
		log.Fatal(err)
	}
}
