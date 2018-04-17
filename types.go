package main

import (
	"context"
	"sync"

	"github.com/nlopes/slack"
	"golang.org/x/sync/errgroup"
)

// Slack configures a Slack client
type Slack struct {
	api *slack.Client
	ctx context.Context
	g   *errgroup.Group

	sync.Mutex
}

// NewSlack creates a new Slack client
func NewSlack(ctx context.Context, api *slack.Client) *Slack {
	var s = new(Slack)
	// set api
	s.api = api

	// set errGroup, and new context to cancel
	s.g, s.ctx = errgroup.WithContext(ctx)

	return s //noop
}

// PostMessage is posting a message to a channel
func (s *Slack) PostMessage(channelID string, message string, params slack.PostMessageParameters) (*errgroup.Group, context.Context) {
	// lock writing
	s.Lock()
	defer s.Unlock()

	s.g.Go(func() error {
		_, _, err := s.api.PostMessageContext(s.ctx, channelID, message, params)

		return err
	})

	return s.g, s.ctx // noop
}
