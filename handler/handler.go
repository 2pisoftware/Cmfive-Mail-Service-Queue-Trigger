package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/log"

	"github.com/2pisoftware/Cmfive-Mail-Service-Queue-Trigger/notifications"
)

// Handle initializes the notifications package and calls the pop method for each sqsEvent.Record.
func Handle(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Create a new structured logger that writes to os.stdout.
	log.Log = log.NewStandardLogger(os.Stdout, log.Fields{
		"tag": "MailService",
	})

	// Initialize the notifications package. Return an error if
	// this fails because we don't want it to be removed from SQS.
	if err := notifications.Initialize(); err != nil {
		log.Error(log.Fields{
			"error": fmt.Sprintf("failed to initialize notifications package: %v", err),
		})
		return fmt.Errorf("failed to initialize notifications package: %w", err)
	}

	// Range over the sqsEvent.Records.
	for _, r := range sqsEvent.Records {
		emailData, err := parse(r)
		if err != nil {
			log.Error(log.Fields{
				"error":        fmt.Sprintf("failed to parse events.SQSMessage into notifications.EmailData: %v", err),
				"sqsMessageID": r.MessageId,
			})

			continue
		}

		messageID, err := notifications.SendEmail(ctx, emailData)
		if err != nil {
			log.Error(log.Fields{
				"error":        fmt.Sprintf("failed to send email: %v", err),
				"sqsMessageID": r.MessageId,
				"to":           emailData.To,
				"cc":           *emailData.CC,
				"bcc":          *emailData.BCC,
				"replyTo":      *emailData.ReplyTo,
				"from":         emailData.From,
				"subject":      emailData.Subject,
			})

			continue
		}

		log.Info(log.Fields{
			"message":      "successfully sent email",
			"sesMessageID": messageID,
		})
	}

	return nil
}

// parse takes the events.SQSMessage, unmarshals it into an notifications.EmailData object.
func parse(message events.SQSMessage) (*notifications.EmailData, error) {
	emailData := &notifications.EmailData{}
	if err := json.Unmarshal([]byte(message.Body), emailData); err != nil {
		return nil, fmt.Errorf("failed to unmashal message body: %w", err)
	}

	// If we haven't been given a ContentType set a default.
	if emailData.ContentType == "" {
		emailData.ContentType = notifications.ContentTypeTextHTML
	}

	return emailData, nil
}
