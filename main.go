package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	e "github.com/axelspringer/vodka-aws/events"
	l "github.com/axelspringer/vodka-aws/lambda"
	log "github.com/sirupsen/logrus"
)

const (
	defaultTimeout     = 60
	defaultEventSource = "aws:sns"
	defaultEnvSSMPath  = "SSM_PATH"
)

// runtime
var (
	ssmPath       string
	ssmEnv        map[string]string
	ssmParameters = []string{"dynamodb-tablename"}
)

// errors
var (
	errNoSSMPath = errors.New("no SSM path configured")
)

// init config
func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// Handler is executed by AWS Lambda in the main function
func Handler(ctx context.Context, event events.SNSEvent) error {
	var err error

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout*time.Second)
	defer cancel()

	// get SSM path from env
	ssmPath, ok := os.LookupEnv(defaultEnvSSMPath)
	if !ok {
		return errNoSSMPath
	}

	// logger
	logger := log.WithFields(log.Fields{
		"ssmPath": ssmPath,
		"event":   event,
	})

	// log
	logger.Info("Configured")

	// create new lambda environment
	λ := l.New(ssmPath)
	if _, err = λ.Store.TestEnv(ssmParameters); err != nil {
		return logError(logger, err)
	}

	// prepare env
	env, err := λ.Store.GetEnv()
	if err != nil {
		return logError(logger, err)
	}

	// prepare aws & dynamo
	sess := session.New()
	db := NewDB(ctx, dynamodb.New(sess), env["dynamodb-tablename"])

	// create signaleer
	signaleer := NewSignaleer(ctx)

	// parse all message records
	for _, record := range event.Records {
		var p = new(e.CodePipelineEvent)
		var hooks []*WebHook

		// filter the events
		if record.EventSource != defaultEventSource {
			logger.Error(err) // log
			continue          // pass along
		}

		// unmarshal the CodePipeline event
		if err := json.Unmarshal([]byte(record.SNS.Message), &p); err != nil {
			logger.Error(err) // log
			continue          // pass along
		}

		// get pipeline
		hooks, err := db.QueryWebHooks(p.Detail, hooks)
		if err != nil {
			logger.Error(err) // log
			continue          // pass along
		}

		// push the found hooks to the signaleer,
		// which then executed the hooks in goroutines
		signaleer.Send(hooks, p.Detail)
	}

	signaleer.Wait() // wait

	return err // noop
}

func logError(logger *log.Entry, err error) error {
	logger.Error(err)
	return err
}

func main() {
	lambda.Start(Handler)
}
