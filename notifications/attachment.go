package notifications

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AttachmentType is a string type alias for the different supported attachment types.
type AttachmentType string

const (
	// AttachmentTypeLocal is used for local attachments. As this is currently
	// deployed as a Lambda function it's primary use is testing.
	AttachmentTypeLocal AttachmentType = "local"
	// AttachmentTypeS3 is used for attachments stored in S3.
	AttachmentTypeS3 AttachmentType = "s3"
	// AttachmentTypeHTTP is used for attachments publicly accessible via HTTP.
	AttachmentTypeHTTP AttachmentType = "http"
)

// Attachment stores a path and type which are used to load it locally.
type Attachment struct {
	Path string         `json:"path"`
	Type AttachmentType `json:"type"`
}

// Load loads an Attachment's data and returns it as a byte slice.
func (a *Attachment) Load(ctx context.Context) ([]byte, error) {
	switch a.Type {
	case AttachmentTypeLocal:
		return ioutil.ReadFile(a.Path)
	case AttachmentTypeS3:
		// Split the path up to get the bucket and key.
		parts := strings.Split(a.Path, "/")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid path for AttachmentTypeS3: %s, path must be in the following format <bucket>/<key>", a.Path)
		}

		buffer := aws.NewWriteAtBuffer(nil)

		input := &s3.GetObjectInput{
			Bucket: aws.String(parts[0]),
			Key:    aws.String(a.Path[len(parts[0])+1:]),
		}

		// Download the Attachment data.
		_, err := s3Downloader.DownloadWithContext(ctx, buffer, input)
		if err != nil {
			return nil, fmt.Errorf("failed to download attachment from S3: %w", err)
		}

		return buffer.Bytes(), nil
	case AttachmentTypeHTTP:
		// Get the Attachment data.
		response, err := http.Get(a.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to perform HTTP Get request on attachment: %w", err)
		}
		defer func() {
			_ = response.Body.Close()
		}()

		return ioutil.ReadAll(response.Body)
	default:
		// Fallthrough to the error case as a.Type is not a supported AttachmentType.
	}

	return nil, fmt.Errorf("unable to load Attachment with unexpected type: %s", a.Type)
}
