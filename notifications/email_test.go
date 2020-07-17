package notifications_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gofor-little/env"

	"github.com/2pisoftware/mail-service-popper/notifications"
)

func TestSendEmail(t *testing.T) {
	if err := env.Load("../.env"); err != nil {
		t.Fatalf("failed to load .env: %v", err)
	}

	splitFunc := func(c rune) bool {
		return c == ','
	}

	to := strings.FieldsFunc(os.Getenv("TEST_TO"), splitFunc)
	cc := strings.FieldsFunc(os.Getenv("TEST_CC"), splitFunc)
	if len(cc) == 0 {
		cc = nil
	}
	bcc := strings.FieldsFunc(os.Getenv("TEST_BCC"), splitFunc)
	if len(bcc) == 0 {
		bcc = nil
	}
	replyTo := strings.FieldsFunc(os.Getenv("TEST_REPLY_TO"), splitFunc)
	if len(replyTo) == 0 {
		replyTo = nil
	}

	bodyText, err := ioutil.ReadFile("test-data/body.txt")
	if err != nil {
		t.Fatalf("failed to read body.txt from file: %v", err)
	}

	bodyHTML, err := ioutil.ReadFile("test-data/body.html")
	if err != nil {
		t.Fatalf("failed to read body.html from file: %v", err)
	}

	attachments := &[]notifications.Attachment{
		{
			Type: notifications.AttachmentTypeLocal,
			Path: "test-data/attachment.docx",
		},
		{
			Type: notifications.AttachmentTypeLocal,
			Path: "test-data/attachment.jpg",
		},
		{
			Type: notifications.AttachmentTypeLocal,
			Path: "test-data/attachment.pdf",
		},
		{
			Type: notifications.AttachmentTypeLocal,
			Path: "test-data/attachment.png",
		},
		{
			Type: notifications.AttachmentTypeLocal,
			Path: "test-data/attachment.txt",
		},
	}

	_, err = notifications.SendEmail(to, &cc, &bcc, &replyTo, os.Getenv("TEST_FROM"), os.Getenv("TEST_SUBJECT"), string(bodyText), notifications.ContentTypeTextPlain, attachments)
	if err != nil {
		t.Fatalf("failed to send email: %v", err)
	}

	_, err = notifications.SendEmail(to, &cc, &bcc, &replyTo, os.Getenv("TEST_FROM"), os.Getenv("TEST_SUBJECT"), string(bodyHTML), notifications.ContentTypeTextHTML, attachments)
	if err != nil {
		t.Fatalf("failed to send email: %v", err)
	}
}
