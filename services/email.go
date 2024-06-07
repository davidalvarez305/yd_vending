package services

import (
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"

	"github.com/davidalvarez305/yd_vending/constants"
)

func SendSMTPEmail(subject, recipient, senderEmail string, data any, templateName string) error {
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	smtpUsername := os.Getenv("GMAIL_USERNAME")
	smtpPassword := os.Getenv("GMAIL_PASSWORD")

	templateFile := constants.PARTIAL_TEMPLATES_DIR + templateName
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("error reading template file: %v", err)
	}

	tmpl, err := template.New("email").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	var htmlMessageBody strings.Builder
	err = tmpl.ExecuteTemplate(&htmlMessageBody, "email", data)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	message := []byte(fmt.Sprintf("From: %s\r\n", senderEmail) +
		fmt.Sprintf("To: %s\r\n", recipient) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		htmlMessageBody.String())

	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, []string{recipient}, message)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}
