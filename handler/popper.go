package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/log"

	"github.com/2pisoftware/mail-service-popper/notifications"
)

// Handle initializes the notifications package and calls the pop method for each sqsEvent.Record.
func Handle(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Create a new structured logger that writes to os.stdout.
	log.Log = log.NewStandardLogger(os.Stdout, log.Fields{
		"tag": "mailService",
	})

	// Initialize the notifications package. Return an error if
	// this fails because we don't want it to be removed from SQS.
	if err := notifications.Initialize(); err != nil {
		log.Error(log.Fields{
			"error":   err.Error(),
			"message": "failed to initialize notifications package",
		})

		return fmt.Errorf("failed to initialize notifications package: %w", err)
	}

	// Range over the sqsEvent.Records.
	for _, r := range sqsEvent.Records {
		// Pop the record.
		messageID, err := pop(ctx, r)
		if err != nil {
			log.Error(log.Fields{
				"error":        err.Error(),
				"message":      "failed to send email",
				"sqsMessageID": r.MessageId,
			})

			// If there is an error attempt to report it.
			if err := reportError(); err != nil {
				log.Error(log.Fields{
					"error":   err.Error(),
					"message": "failed to report error",
				})
			}

			continue
		}

		log.Info(log.Fields{
			"message":      "successfully sent email",
			"sesMessageID": messageID,
		})
	}

	return nil
}

// pop takes the events.SQSMessage, unmarshals it into a notifications.EmailData
// structure and passes it to notifications.SendEmail.
func pop(ctx context.Context, message events.SQSMessage) (string, error) {
	// Unmarshal the message body into notifications.EmailData.
	emailData := &notifications.EmailData{}
	if err := json.Unmarshal([]byte(message.Body), emailData); err != nil {
		return "", fmt.Errorf("failed to unmashal message body: %w", err)
	}

	// If we haven't been given a ContentType set a default.
	if emailData.ContentType == "" {
		emailData.ContentType = notifications.ContentTypeTextHTML
	}

	// Send the email.
	messageID, err := notifications.SendEmail(ctx, emailData)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return messageID, nil
}

// reportError will send a message to an endpoint informing of the failure to send an email.
func reportError() error {
	// TODO: Implement in a future phase.
	return nil
}
