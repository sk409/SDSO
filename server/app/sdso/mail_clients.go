package main

import (
	"net/smtp"
	"os"
	"strings"
)

type mailClient interface {
	send(to, subject, body string, html bool) error
}

func newMailClient(email string) (mailClient, error) {
	if strings.HasSuffix(email, "@gmail.com") {
		return &gmailClient{}, nil
	}
	return nil, errUnkownDomain
}

type gmailClient struct {
}

func (g *gmailClient) send(to, subject, body string, html bool) error {
	gmailUser := os.Getenv("GMAIL_USER")
	gmailAppPassword := os.Getenv("GMAIL_APP_PASSWORD")
	gmailHost := "smtp.gmail.com"
	gmailServer := gmailHost + ":587"
	auth := smtp.PlainAuth("", gmailUser, gmailAppPassword, gmailHost)
	mime := ""
	if html {
		mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	}
	subject = "Subject: " + subject + "\n"
	msg := []byte(subject + mime + body)
	err := smtp.SendMail(gmailServer, auth, "", []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}
