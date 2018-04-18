package main

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	e "github.com/axelspringer/vodka-aws/events"
)

// DB represents a DynamoDB
type DB struct {
	ctx       context.Context
	db        *dynamodb.DynamoDB
	tableName string
}

// NewDB is returning a new DB
func NewDB(ctx context.Context, db *dynamodb.DynamoDB, tableName string) *DB {
	return &DB{ctx, db, tableName}
}

// GetSlack is getting a pipeline from the DynamoDB table
func (d *DB) GetSlack(event e.CodePipelineEventDetails, slack *Slack) (*Slack, error) {
	var err error

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Pipeline": {
				S: aws.String(event.Pipeline),
			},
		},
		TableName: aws.String(d.tableName),
	}

	output, err := d.db.GetItemWithContext(d.ctx, input)
	if err != nil {
		return slack, err
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, slack)

	return slack, err // noop
}
