package main

import (
	"fmt"
	"net/http"

	e "github.com/axelspringer/vodka-aws/events"
	"github.com/axelspringer/vodka-aws/vodka"
	"github.com/dghubble/sling"
)

const (
	// SlackDefaultColor represents the default color
	SlackDefaultColor = "#d3d3d3"
	// SlackWarningColor represents the warning color
	SlackWarningColor = "warning"
	// SlackFailureColor represents the failure color
	SlackFailureColor = "danger"
	// SlackSuccessColor represents the success color
	SlackSuccessColor = "good"
)

// SlackColors represents the mapping of status colors to
var SlackColors = map[string]string{
	e.CodePipelineSucceeded:  SlackSuccessColor,
	e.CodePipelineFailed:     SlackFailureColor,
	e.CodePipelineResumed:    SlackWarningColor,
	e.CodePipelineSuperseded: SlackDefaultColor,
	e.CodePipelineCanceled:   SlackWarningColor,
}

// WebHookPayload represents the general interface to a webhook payload
type WebHookPayload interface {
	Post(webhookURL string) (bool, *http.Response, error)
}

// SlackField represents the field to a SlackAttachment
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// SlackAttachment represents a Slack message attachment
type SlackAttachment struct {
	Fallback   *string       `json:"fallback"`
	Color      *string       `json:"color"`
	PreText    *string       `json:"pretext"`
	AuthorName *string       `json:"author_name"`
	AuthorLink *string       `json:"author_link"`
	AuthorIcon *string       `json:"author_icon"`
	Title      *string       `json:"title"`
	TitleLink  *string       `json:"title_link"`
	Text       *string       `json:"text"`
	ImageURL   *string       `json:"image_url"`
	Fields     []*SlackField `json:"fields"`
	Footer     *string       `json:"footer"`
	FooterIcon *string       `json:"footer_icon"`
	Timestamp  *int64        `json:"ts"`
	MarkdownIn *[]string     `json:"mrkdwn_in"`
}

// SlackPayload represents the payload of a Slack message
type SlackPayload struct {
	Parse       string            `json:"parse,omitempty"`
	Username    string            `json:"username,omitempty"`
	IconURL     string            `json:"icon_url,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	Channel     string            `json:"channel,omitempty"`
	Text        string            `json:"text,omitempty"`
	LinkNames   string            `json:"link_names,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
	UnfurlLinks bool              `json:"unfurl_links,omitempty"`
	UnfurlMedia bool              `json:"unfurl_media,omitempty"`
	Markdown    bool              `json:"mrkdwn,omitempty"`
}

// NewSlackPayload is returning a new SlackPayload
// channel: is the Slack channel to use
// bot: is the name of the bot, mapped to the Slack username displayed
// event: is the CodePipeline Event to display
func NewSlackPayload(channel string, bot string, event e.CodePipelineEventDetail) *SlackPayload {
	attachment := SlackAttachment{
		Color: vodka.String(SlackColors[event.State]),
	}
	attachment.AddField(SlackField{Title: "Pipeline", Value: event.Pipeline})
	attachment.AddField(SlackField{Title: "State", Value: event.State})
	attachment.AddField(SlackField{Title: "Stage", Value: event.Stage})

	return &SlackPayload{
		Channel:     channel,
		Username:    bot,
		Text:        fmt.Sprintf("%s is %s", event.Pipeline, event.State), // todo: golang template
		Attachments: []SlackAttachment{attachment},
	}
}

// AddField adds a field to a SlackAttachment
func (a *SlackAttachment) AddField(field SlackField) *SlackAttachment {
	a.Fields = append(a.Fields, &field)
	return a
}

// Post is posting the payload to a Slack webhook
func (s *SlackPayload) Post(webhookURL string) (bool, *http.Response, error) {
	success := new(interface{})
	resp, err := sling.New().Post(webhookURL).BodyJSON(s).ReceiveSuccess(success)

	return resp.StatusCode == 200, resp, err
}

// AddAttachment is adding an new attachment to a SlackPayload
func (s *SlackPayload) AddAttachment(attachment SlackAttachment) *SlackPayload {
	s.Attachments = append(s.Attachments, attachment)
	return s
}
