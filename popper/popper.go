package popper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/2pisoftware/mail-service-popper/notifications"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofor-little/log"
)

func HandlePop(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Log = log.NewStandardLogger(os.Stdout, log.Fields{
		"tag": "mailService",
	})

	for _, r := range sqsEvent.Records {
		messageID, err := pop(r)
		if err != nil {
			log.Error(log.Fields{
				"error":        err.Error(),
				"message":      "failed to send email",
				"sqsMessageID": r.MessageId,
			})

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

func pop(message events.SQSMessage) (string, error) {
	emailData := &EmailData{}
	if err := json.Unmarshal([]byte(message.Body), emailData); err != nil {
		return "", fmt.Errorf("failed to unmashal message body: %w", err)
	}

	// TODO: Remove once it's coming from Cmfive.
	if emailData.ContentType == "" {
		emailData.ContentType = notifications.ContentTypeTextHTML
	}

	messageID, err := notifications.SendEmail(emailData.To, emailData.CC, emailData.BCC, emailData.ReplyTo, emailData.From, emailData.Subject, emailData.Body, emailData.ContentType, emailData.Attachments)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return messageID, nil
}

// TODO: Handle error by sending details to Cmfive endpoint.
func reportError() error {
	return nil
}
