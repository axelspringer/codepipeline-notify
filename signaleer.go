package main

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nlopes/slack"

	e "github.com/axelspringer/vodka-aws/events"
)

// Signal defines a signal to be processed
type Signal struct {
	TableName string
	Detail    e.CodePipelineEventDetails
}

// Signaleer contains the service to signal Slack
type Signaleer struct {
	ctx     context.Context
	db      *dynamodb.DynamoDB
	session *session.Session

	wg sync.WaitGroup
	sync.Mutex
}

// NewSignaleer creates a new Signaleer to be used to signal Slack channels about pipelines events
func NewSignaleer(ctx context.Context, session *session.Session, db *dynamodb.DynamoDB) *Signaleer {
	return &Signaleer{ctx: ctx, db: db, session: session}
}

// PostMessage is posting a message to a Slack Channel
func (s *Signaleer) PostMessage(api *slack.Client) {
	s.Lock() // make it safe
	defer s.Unlock()

	wg.Add(1)
}

// GetPipeline is getting a pipeline from the DynamoDB table
func (s *Signaleer) GetPipeline(signal *Signal) (*dynamodb.GetItemOutput, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Pipeline": {
				S: aws.String(signal.Detail.Pipeline),
			},
		},
		TableName: aws.String(signal.TableName),
	}

	return s.db.GetItemWithContext(s.ctx, input)
}

// Wait is using the WaitGroup to wait for all message to execute
func (s *Signaleer) Wait() {
	wg.Wait()
}
