package main_test

import (
	"context"
	"errors"
	"os"
	"time"

	main "github.com/axelspringer/codepipeline-notify"
	event "github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/snsevt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

const (
	defaultEnvSSMPath = "SSM_PATH"
)

// events
var (
	eventSuccess = event.Event{
		Records: []*event.EventRecord{&event.EventRecord{
			SNS: &event.Record{},
		}},
	}
)

// SSM paths
var (
	pathSuccess = "/mock"
)

// errors
var (
	errNoSSMPath = errors.New("no SSM path configured")
)

// testing Events on the
func testingEvents(e event.Event, ssmPath *string, err error) {
	// setup end
	if ssmPath != nil {
		os.Setenv(defaultEnvSSMPath, *ssmPath)
		defer os.Unsetenv(defaultEnvSSMPath)
	}
	// create new ctx, and run handler
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	λ := main.Handler(ctx, e)

	// handle the non-error
	if err == nil {
		Ω(λ).Should(Succeed())
	}

	// handle the error case
	if err != nil {
		Ω(λ).Should(MatchError(err))
	}
}

var _ = Describe("Handler", func() {
	DescribeTable("events",
		testingEvents,
		Entry("success", eventSuccess, &pathSuccess, nil),
		Entry("failure", eventSuccess, nil, errNoSSMPath),
	)
})
