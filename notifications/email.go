package notifications

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

// ContentType is a string type alias for the different supported content types.
type ContentType string

const (
	// ContentTypeTextPlain is used for text/plain emails.
	ContentTypeTextPlain ContentType = "text/plain"
	// ContentTypeTextHTML is used for text/html emails.
	ContentTypeTextHTML ContentType = "text/html"
)

// SendEmail loads the Attachments, builds an Email and sends it via SES.
func SendEmail(ctx context.Context, data *EmailData) (string, error) {
	// Check that the package has been initialized.
	if sesClient == nil || s3Downloader == nil {
		return "", errors.New("notifications.Initialize has not been called yet")
	}

	// Build the destination emails.
	var destinations []string
	destinations = append(destinations, data.To...)

	if data.CC != nil {
		destinations = append(destinations, *data.CC...)
	}
	if data.BCC != nil {
		destinations = append(destinations, *data.BCC...)
	}

	// Set the headers and body.
	message := gomail.NewMessage()
	message.SetHeader("To", destinations...)
	message.SetHeader("Reply-To", *data.ReplyTo...)
	message.SetHeader("From", data.From)
	message.SetHeader("Subject", data.Subject)
	message.SetBody(string(data.ContentType), data.Body)

	// Load and attach the Attachments.
	for _, a := range *data.Attachments {
		path, err := a.Load(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to load Attachment: %w", err)
		}

		message.Attach(path)
	}

	// Remove any attachments on disc once we're done with them.
	defer func() {
		// Don't delete attachments when we're developing locally.
		if os.Getenv("ENVIRONMENT") == "development" {
			return
		}

		for _, a := range *data.Attachments {
			_ = os.Remove(a.Path)
		}
	}()

	// Write the email to a buffer.
	var buf bytes.Buffer
	_, err := message.WriteTo(&buf)
	if err != nil {
		return "", err
	}

	// Create and validate raw email input.
	input := &ses.SendRawEmailInput{
		Destinations: aws.StringSlice(destinations),
		FromArn:      aws.String(os.Getenv("AWS_FROM_ARN")),
		RawMessage: &ses.RawMessage{
			Data: buf.Bytes(),
		},
	}

	if err := input.Validate(); err != nil {
		return "", fmt.Errorf("failed to validate ses.SendRawEmailInput: %w", err)
	}

	// Send email.
	output, err := sesClient.SendRawEmail(input)
	if err != nil {
		return "", fmt.Errorf("failed to send raw email: %w", err)
	}

	return *output.MessageId, nil
}
