package notifications

import (
	"fmt"
	"os"

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
		return fmt.Errorf("failed to create new session.Session: %w", err)
	}

	sesClient = ses.New(sess)
	s3Downloader = s3manager.NewDownloader(sess)

	return nil
}
