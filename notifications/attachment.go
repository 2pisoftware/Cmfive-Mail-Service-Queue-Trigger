package notifications

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofor-little/rand"
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

// Load loads an Attachment's data and stores is locally in /tmp to be used.
func (a *Attachment) Load(ctx context.Context) (string, error) {
	switch a.Type {
	case AttachmentTypeLocal:
		// Just return the local file path.
		return a.Path, nil
	case AttachmentTypeS3:
		// Split the path up to get the bucket and key.
		parts := strings.Split(a.Path, "/")
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid path for AttachmentTypeS3: %s, path must be in the following format <bucket>/<key>", a.Path)
		}

		buffer := aws.NewWriteAtBuffer(nil)

		input := &s3.GetObjectInput{
			Bucket: aws.String(parts[0]),
			Key:    aws.String(a.Path[len(parts[0])+1:]),
		}

		// Download the Attachment data.
		_, err := s3Downloader.DownloadWithContext(ctx, buffer, input)
		if err != nil {
			return "", fmt.Errorf("failed to download attachment from S3: %w", err)
		}

		// Find a free file path for the Attachment file.
		freePath, err := a.getAvailablePath()
		if err != nil {
			return "", fmt.Errorf("failed to get free path: %w", err)
		}

		// Create the Attachment file in the /tmp directory.
		file, err := os.Create(freePath)
		if err != nil {
			return "", fmt.Errorf("failed to create attachment file: %w", err)
		}
		defer func() {
			_ = file.Close()
		}()

		// Write the Attachment data to file.
		_, err = file.Write(buffer.Bytes())
		if err != nil {
			return "", fmt.Errorf("failed to write buffer to file: %w", err)
		}

		return freePath, nil
	case AttachmentTypeHTTP:
		// Get the Attachment data.
		response, err := http.Get(a.Path)
		if err != nil {
			return "", fmt.Errorf("failed to perform HTTP Get request on attachment: %w", err)
		}
		defer func() {
			_ = response.Body.Close()
		}()

		// Find a free file path for the Attachment file.
		freePath, err := a.getAvailablePath()
		if err != nil {
			return "", fmt.Errorf("failed to get free path: %w", err)
		}

		// Create the Attachment file in the /tmp directory.
		file, err := os.Create(freePath)
		if err != nil {
			return "", fmt.Errorf("failed to create attachment file: %w", err)
		}
		defer func() {
			_ = file.Close()
		}()

		// Copy the Attachment data to file.
		_, err = io.Copy(file, response.Body)
		if err != nil {
			return "", fmt.Errorf("failed to copy HTTP response body to file: %w", err)
		}

		return freePath, nil
	default:
		// Fallthrough to the error case as a.Type is not a supported AttachmentType.
	}

	return "", fmt.Errorf("unable to load Attachment with unexpected type: %s", a.Type)
}

// getAvailablePath will ensure the path that the file will be stored in doesn't
// collide with another file.
func (a *Attachment) getAvailablePath() (string, error) {
	// Get the base or file name from the path.
	base := path.Base(a.Path)
	// If we don't get a valid file name create one.
	if base == "." || base == "/" {
		var err error
		base, err = rand.GenerateCryptoString(16)
		if err != nil {
			return "", err
		}
	}

	freePath := fmt.Sprintf("/tmp/%s", base)

	// Check if the path is available, if so return that path.
	if _, err := os.Stat(freePath); err != nil {
		if os.IsNotExist(err) {
			return freePath, nil
		}

		return "", err
	}

	index := 0
	isFreePath := false

	// If the path isn't available append the index to the filename (without the
	// extension) until an available path is found. For example, file_1.txt, file_2.txt.
	for !isFreePath {
		extension := filepath.Ext(freePath)
		name := freePath[0 : len(freePath)-len(extension)]

		fp := fmt.Sprintf("%s_%d%s", name, index, extension)

		if _, err := os.Stat(fp); err != nil {
			if os.IsNotExist(err) {
				isFreePath = true
				freePath = fp
				continue
			}

			return "", err
		}

		index++
	}

	return freePath, nil
}
