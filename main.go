package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	l "github.com/axelspringer/vodka-aws/lambda"
	event "github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/snsevt"
	log "github.com/sirupsen/logrus"
)

const (
	defaultEnvSSMPath = "SSM_PATH"
)

// runtime
var (
	ssmPath       string
	ssmEnv        map[string]string
	ssmParameters = []string{"hook", "channel"}
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
func Handler(ctx context.Context, event event.Event) error {
	var err error

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
	_, err = λ.Store.GetEnv()
	if err != nil {
		return logError(logger, err)
	}

	return err // noop
}

func logError(logger *log.Entry, err error) error {
	logger.Error(err)
	return err
}

func main() {
	lambda.Start(Handler)
}
