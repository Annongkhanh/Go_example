package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress = "smtp.gmail.com"
	smtpAuthPort    = 587
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	e.HTML = []byte(content)

	for _, attachFile := range attachFiles {
		_, err := e.AttachFile(attachFile)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", attachFile, err)
		}
	}
	auth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	smtpAuthServer := fmt.Sprintf("%s:%d", smtpAuthAddress, smtpAuthPort)
	return e.Send(smtpAuthServer, auth)
}
