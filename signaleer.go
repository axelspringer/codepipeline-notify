package main

import (
	"context"
	"fmt"
	"sync"

	e "github.com/axelspringer/vodka-aws/events"
)

// Signaleer contains the service to signal Slack
type Signaleer struct {
	ctx context.Context

	wg sync.WaitGroup
	sync.Mutex
}

// Signal contains the interface to a signal
type Signal interface {
	Send(ctx context.Context) error
}

// NewSignaleer creates a new Signaleer to be used to signal Slack channels about pipelines events
func NewSignaleer(ctx context.Context) *Signaleer {
	return &Signaleer{ctx: ctx}
}

// Send is posting a message to a Slack Channel
func (s *Signaleer) Send(webHooks []*WebHook, event e.CodePipelineEventDetails) {
	s.Lock() // safe
	defer s.Unlock()

	wg.Add(1) // new routine

	for _, hook := range webHooks {
		go func(event e.CodePipelineEventDetails) {
			err := hook.Send(s.ctx, event)

			if err != nil {
				fmt.Println(err)
			}

			s.Lock() // safe
			defer s.Unlock()

			wg.Done()
		}(event)
	}
}

// Wait is using the WaitGroup to wait for all message to execute
func (s *Signaleer) Wait() {
	wg.Wait()
}
