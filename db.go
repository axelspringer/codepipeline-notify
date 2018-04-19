package main

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

// QueryWebHooks is getting a pipeline from the DynamoDB table
func (d *DB) QueryWebHooks(event e.CodePipelineEventDetail, hooks []*WebHook) ([]*WebHook, error) {
	var err error

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(event.Pipeline),
			},
		},
		KeyConditionExpression: aws.String("Pipeline = :v1"),
		TableName:              aws.String(d.tableName),
	}

	query, err := d.db.QueryWithContext(d.ctx, input)
	if err != nil {
		return hooks, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(query.Items, &hooks)

	return hooks, err // noop
}
