package scheduler

import (
	"context"
	"fmt"

	"github.com/Justified02/abm/internal/fetcher"
)

type fetchResult struct {
	data []byte
	err  error
}

type Scheduler struct {
	stripe *fetcher.StripeClient
}

func NewScheduler(s *fetcher.StripeClient) *Scheduler {
	newScheduler := &Scheduler{
		stripe: s,
	}

	return newScheduler
}

func (s *Scheduler) runCollection(ctx context.Context) {
	fetchedResult := make(chan fetchResult, 1)

	go func() {
		data, err := s.stripe.Fetch(ctx)
		fetchedResult <- fetchResult{data: data, err: err}
	}()

	result := <-fetchedResult

	if result.err != nil {
		fmt.Println("error fetching data:", result.err)
		return
	}

	fmt.Println(string(result.data))
}
