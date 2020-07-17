package popper_test

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/2pisoftware/mail-service-popper/popper"
)

func TestPop(t *testing.T) {
	if err := popper.HandlePop(context.Background(), events.SQSEvent{}); err != nil {
		t.Fatal(err)
	}
}
