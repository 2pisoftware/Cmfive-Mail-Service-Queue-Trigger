package notifications

import (
	"fmt"

	"github.com/gofor-little/env"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	sesClient    *ses.SES
	s3Downloader *s3manager.Downloader
)

// Initialize initializes this package by creating a new SES client and S3 downloader.
func Initialize() error {
	var sess *session.Session
	var err error

	if env.Get("ENVIRONMENT", "production") == "development" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(env.Get("AWS_REGION", "ap-southeast-2")),
			},
			Profile: env.Get("AWS_PROFILE", "default"),
		})
	} else {
		sess, err = session.NewSession()
	}

	if err != nil {
		return fmt.Errorf("failed to create new session.Session: %w", err)
	}

	sesClient = ses.New(sess)
	s3Downloader = s3manager.NewDownloader(sess)

	return nil
}
