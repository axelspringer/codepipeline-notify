package main

import (
	"context"
	"fmt"

	e "github.com/axelspringer/vodka-aws/events"
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

// sendPayload is sending a WebHookPayload
func (s *WebHook) sendPayload(url string, payload WebHookPayload) error {
	var err error

	resp, err := payload.Send(url)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error sending payload. Status: %v", resp.Status)
	}

	return err
}

// Send is sending a signal
func (s *WebHook) Send(ctx context.Context, event e.CodePipelineEventDetails) error {
	var err error

	switch webhookType := s.Type; webhookType {
	case WebHookTypeSlack:
		slack := NewSlackPayload(s.Channel, s.Bot, event)
		err = s.sendPayload(s.URL, slack)
	default:
	}

	return err
}
