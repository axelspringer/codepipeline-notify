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

	log *log.Entry
}

// postPayload is posting a WebHookPayload
func (w *WebHook) postPayload(webHookURL string, payload WebHookPayload) error {
	var err error

	success, _, err := payload.Post(webHookURL) // omit response for now
	if !success || err != nil {
		return err
	}

	return err // noop
}

// Send is sending a signal
func (w *WebHook) Send(ctx context.Context, event e.CodePipelineEventDetail) error {
	var err error

	// setup logger
	w.log = log.WithFields(log.Fields{
		"type":     w.Type,
		"pipeline": w.Pipeline,
		"channel":  w.Channel,
		"url":      w.URL,
	})

	// log
	w.log.Info("Executing WebHook")

	// configure the webhook
	switch webhookType := w.Type; webhookType {
	case WebHookTypeSlack:
		slack := NewSlackPayload(w.Channel, w.Bot, event)
		err = w.postPayload(w.URL, slack)
		if err != nil {
			w.log.WithError(err)
		}
	default:
	}

	return err // noop
}
