package notifications

// EmailData is the expected result once the JSON data has been
// decoded from SQS.
type EmailData struct {
	To          []string      `json:"to"`
	CC          *[]string     `json:"cc,omitempty"`
	BCC         *[]string     `json:"bcc,omitempty"`
	ReplyTo     *[]string     `json:"reply_to,omitempty"`
	From        string        `json:"from"`
	Subject     string        `json:"subject"`
	Body        string        `json:"body"`
	ContentType ContentType   `json:"content_type"`
	Attachments *[]Attachment `json:"attachments,omitempty"`
}
