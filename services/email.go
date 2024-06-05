package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendSMTPEmail(subject, body, recipient string) error {
	// SMTP server configuration
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	smtpUsername := os.Getenv("GMAIL_USERNAME")
	smtpPassword := os.Getenv("GMAIL_PASSWORD")
	senderEmail := os.Getenv("GMAIL_EMAIL")

	htmlMessage := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>HTML Email</title>
		</head>
		<body>
			<h1>Hello!</h1>
			<p>This is a test HTML email.</p>
		</body>
		</html>
	`

	// Authentication
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	// Compose MIME message
	message := []byte(fmt.Sprintf("From: %s\r\n", senderEmail) +
		fmt.Sprintf("To: %s\r\n", "recipient@example.com") +
		"Subject: Test HTML Email\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		htmlMessage)

	// Send email
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, []string{recipient}, message)
	if err != nil {
		return err
	}

	return nil
}
