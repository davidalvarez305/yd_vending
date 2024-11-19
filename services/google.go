package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func refreshAuthToken(config *oauth2.Config) (oauth2.Token, error) {
	var token oauth2.Token

	refreshToken := constants.GoogleRefreshToken
	client := &http.Client{}

	url := config.Endpoint.TokenURL
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Request failed: ", err)
		return token, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("client_id", config.ClientID)
	q.Add("client_secret", config.ClientSecret)
	q.Add("refresh_token", refreshToken)
	q.Add("grant_type", "refresh_token")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while getting auth token: ", err)
		return token, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("STATUS CODE: %+v\n", resp.Status)
		return token, errors.New("request failed")
	}

	json.NewDecoder(resp.Body).Decode(&token)

	return token, nil
}

func initializeGoogleClient(scope string) (*http.Client, error) {
	var client *http.Client

	b, err := os.ReadFile(constants.GoogleJSONPath)
	if err != nil {
		return client, err
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		return client, err
	}

	token, err := refreshAuthToken(config)

	if err != nil {
		return nil, err
	}

	return config.Client(context.Background(), &token), nil
}

func SendGmail(recipients []string, subject, sender, body string) error {
	client, err := initializeGoogleClient(gmail.GmailSendScope)
	if err != nil {
		fmt.Printf("Unable to initialize Gmail client: %v", err)
		return err
	}

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Gmail client: %v", err)
		return err
	}

	user := "me"

	var message gmail.Message

	emailContent := fmt.Sprintf("To: %s\r\nSubject: %s\r\nReply-To:%s\r\n%s", strings.Join(recipients, ", "), subject, sender, body)
	message.Raw = base64.URLEncoding.EncodeToString([]byte(emailContent))

	_, err = srv.Users.Messages.Send(user, &message).Do()
	if err != nil {
		fmt.Printf("Unable to send email: %v", err)
		return err
	}

	return nil
}

func SendGmailWithAttachment(recipients []string, subject, sender, body, attachmentPath string) error {
	client, err := initializeGoogleClient(gmail.GmailSendScope)
	if err != nil {
		fmt.Printf("Unable to initialize Gmail client: %v", err)
		return err
	}

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Gmail client: %v", err)
		return err
	}

	var msgBuilder strings.Builder

	// Email headers
	msgBuilder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(recipients, ", ")))
	msgBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msgBuilder.WriteString(fmt.Sprintf("Reply-To: %s\r\n", sender))
	msgBuilder.WriteString("MIME-Version: 1.0\r\n")
	msgBuilder.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", constants.EmailMIMEBoundary))
	msgBuilder.WriteString("\r\n")

	// Email body part
	msgBuilder.WriteString(fmt.Sprintf("--%s\r\n", constants.EmailMIMEBoundary))
	msgBuilder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	msgBuilder.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	msgBuilder.WriteString(body)
	msgBuilder.WriteString("\r\n\r\n")

	if attachmentPath != "" {
		fileBytes, err := os.ReadFile(attachmentPath)
		if err != nil {
			return fmt.Errorf("unable to read attachment file: %v", err)
		}

		msgBuilder.WriteString(fmt.Sprintf("--%s\r\n", constants.EmailMIMEBoundary))
		msgBuilder.WriteString("Content-Type: application/octet-stream\r\n")
		msgBuilder.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%q\r\n", attachmentPath))
		msgBuilder.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")

		attachmentEncoded := base64.StdEncoding.EncodeToString(fileBytes)
		msgBuilder.WriteString(attachmentEncoded)
		msgBuilder.WriteString("\r\n\r\n")
	}

	msgBuilder.WriteString(fmt.Sprintf("--%s--", constants.EmailMIMEBoundary))

	message := gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(msgBuilder.String())),
	}

	user := "me"
	_, err = srv.Users.Messages.Send(user, &message).Do()
	if err != nil {
		fmt.Printf("Unable to send email: %v", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}

func ScheduleGoogleCalendarEvent(eventTitle, description, location string, startTime, endTime time.Time, attendees []string) (string, error) {
	var link string
	client, err := initializeGoogleClient(calendar.CalendarScope)
	if err != nil {
		fmt.Printf("Unable to initialize Google Calendar client: %v", err)
		return link, err
	}

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to create Calendar service: %v", err)
		return link, err
	}

	event := &calendar.Event{
		Summary:     eventTitle,
		Location:    location,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: constants.TimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: constants.TimeZone,
		},
		Attendees: []*calendar.EventAttendee{},

		Reminders: &calendar.EventReminders{
			UseDefault: true,
		},
	}

	for _, email := range attendees {
		event.Attendees = append(event.Attendees, &calendar.EventAttendee{
			Email: email,
		})
	}

	createdEvent, err := srv.Events.Insert("primary", event).Do()
	if err != nil {
		fmt.Printf("Unable to create event: %v", err)
		return link, err
	}

	link = createdEvent.HtmlLink

	return link, err
}
