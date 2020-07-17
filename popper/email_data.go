package popper

import "github.com/2pisoftware/mail-service-popper/notifications"

type EmailData struct {
	To          []string                    `json:"to"`
	CC          *[]string                   `json:"cc,omitempty"`
	BCC         *[]string                   `json:"bcc,omitempty"`
	ReplyTo     *[]string                   `json:"reply_to,omitempty"`
	From        string                      `json:"from"`
	Subject     string                      `json:"subject"`
	Body        string                      `json:"body"`
	ContentType notifications.ContentType   `json:"content_type"`
	Attachments *[]notifications.Attachment `json:"attachments,omitempty"`
}
