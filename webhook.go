package main

import (
	"context"

	e "github.com/axelspringer/vodka-aws/events"
	log "github.com/sirupsen/logrus"
)

const (
	// WebHookTypeSlack represents a Slack WebHook
	WebHookTypeSlack = "Slack"
)

// WebHook contains a WebHook config
type WebHook struct {
	Pipeline string
	Channel  string
	Bot      string
	URL      string
	Type     string
}

// postPayload is posting a WebHookPayload
func (w *WebHook) postPayload(webHookURL string, payload WebHookPayload) error {
	var err error

	// request log
	logger := log.WithFields(log.Fields{
		"webHookType": w.Type,
	})

	success, _, err := payload.Post(webHookURL) // omit response for now
	if !success || err != nil {
		logger.WithError(err)
		return err
	}

	return err
}

// Send is sending a signal
func (w *WebHook) Send(ctx context.Context, event e.CodePipelineEventDetail) error {
	var err error

	switch webhookType := w.Type; webhookType {
	case WebHookTypeSlack:
		slack := NewSlackPayload(w.Channel, w.Bot, event)
		err = w.postPayload(w.URL, slack)
	default:
	}

	return err
}
