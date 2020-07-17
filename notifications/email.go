package notifications

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

type ContentType string

const (
	ContentTypeTextPlain ContentType = "text/plain"
	ContentTypeTextHTML  ContentType = "text/html"
)

func SendEmail(to []string, cc *[]string, bcc *[]string, replyTo *[]string, from string, subject string, body string, contentType ContentType, attachments *[]Attachment) (string, error) {
	client, err := newSESClient()
	if err != nil {
		return "", fmt.Errorf("failed to create new SES client: %w", err)
	}

	var destinations []string
	destinations = append(destinations, to...)

	if cc != nil {
		destinations = append(destinations, *cc...)
	}
	if bcc != nil {
		destinations = append(destinations, *bcc...)
	}

	message := gomail.NewMessage()
	message.SetHeader("To", destinations...)
	message.SetHeader("Reply-To", *replyTo...)
	message.SetHeader("From", from)
	message.SetHeader("Subject", subject)
	message.SetBody(string(contentType), body)

	for _, a := range *attachments {
		path, err := a.Load()
		if err != nil {
			return "", fmt.Errorf("failed to load Attachment: %w", err)
		}

		message.Attach(path)
	}

	var data bytes.Buffer
	_, err = message.WriteTo(&data)
	if err != nil {
		return "", err
	}

	input := &ses.SendRawEmailInput{
		Destinations: aws.StringSlice(destinations),
		FromArn:      aws.String(os.Getenv("AWS_FROM_ARN")),
		RawMessage: &ses.RawMessage{
			Data: data.Bytes(),
		},
	}

	if err := input.Validate(); err != nil {
		return "", fmt.Errorf("failed to validate ses.SendRawEmailInput: %w", err)
	}

	output, err := client.SendRawEmail(input)
	if err != nil {
		return "", fmt.Errorf("failed to send raw email: %w", err)
	}

	return *output.MessageId, nil
}

func newSESClient() (*ses.SES, error) {
	var sess *session.Session
	environment := os.Getenv("GO_ENV")
	var err error

	if environment == "development" || environment == "" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(os.Getenv("AWS_REGION")),
			},
			Profile: os.Getenv("AWS_PROFILE"),
		})
	} else {
		sess, err = session.NewSession()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create new session.Session: %w", err)
	}

	return ses.New(sess), nil
}
