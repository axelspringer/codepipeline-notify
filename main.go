package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	event "github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/codepipelineevt"
	log "github.com/sirupsen/logrus"
)

const (
	defaultEnvSSMPath = "SSM_PATH"
)

// runtime
var (
	ssmPath string
)

// errors
var (
	errNoSSMPath = errors.New("no SSM path configured")
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(ctx context.Context, event event.Event) error {
	var err error

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

	return err // noop
}

func main() {
	lambda.Start(Handler)
}
