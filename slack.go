package main

import (
	"context"
	"fmt"

	"github.com/nlopes/slack"
)

// Slack contains a Slack config
type Slack struct {
	Pipeline string
	Channel  string
	Bot      string
	Token    string

	api *slack.Client
}

// Send is sending a signal
func (s *Slack) Send(ctx context.Context) error {
	var err error

	s.newAPI() // create new api, if not exists

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext: "some pretext",
		Text:    "some text",
		// Uncomment the following part to send a field too
		/*
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "a",
					Value: "no",
				},
			},
		*/
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err = s.api.PostMessageContext(ctx, s.Channel, "Some text", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}

	return nil
}

// newAPI is setting a new API
func (s *Slack) newAPI() {
	if s.api == nil {
		s.api = slack.New(s.Token)
	}
}
