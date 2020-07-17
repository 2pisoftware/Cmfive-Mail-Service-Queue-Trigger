package notifications

import "fmt"

type AttachmentType string

const (
	AttachmentTypeLocal AttachmentType = "local"
	AttachmentTypeS3    AttachmentType = "s3"
)

type Attachment struct {
	Path string         `json:"path"`
	Type AttachmentType `json:"type"`
}

func (a *Attachment) Load() (string, error) {
	switch a.Type {
	case AttachmentTypeLocal:
		return a.Path, nil
	default:
	}

	return "", fmt.Errorf("unable to load Attachment with unexpected type: %s", a.Type)
}
