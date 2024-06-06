package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/davidalvarez305/budgeting/database"
	"github.com/davidalvarez305/budgeting/middleware"
	"github.com/davidalvarez305/budgeting/models"
	"github.com/davidalvarez305/budgeting/services"
)

func PhoneServiceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/call/inbound":
			handleIncomingCall(w, r)
		case "/sms/inbound":
			handleIncomingSMS(w, r)
		case "/sms/outbound":
			handleOutboundSMS(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleIncomingCall(w http.ResponseWriter, r *http.Request) {

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	from := r.FormValue("From")
	to := r.FormValue("To")

	forwardTo := os.Getenv("DAVID_TWILIO_PHONE_NUMBER")

	fmt.Printf("Incoming call from: %s", from)
	fmt.Printf("Incoming call to: %s", to)

	twiML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Dial>%s</Dial>
	</Response>`, forwardTo)

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(twiML))
}

func handleIncomingSMS(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	userId, err := database.GetUserIDFromPhoneNumber(r.FormValue("From"))
	if err != nil {
		http.Error(w, "Failed to get User ID.", http.StatusBadRequest)
		return
	}

	sms := models.TextMessage{
		MessageSID: r.FormValue("MessageSid"),
		UserID:     userId,
		FromNumber: r.FormValue("From"),
		ToNumber:   r.FormValue("To"),
		Body:       r.FormValue("Body"),
		Status:     "received",
		CreatedAt:  time.Now(),
		IsInbound:  true,
	}

	// Save the SMS to the database
	if err := database.SaveSMS(sms); err != nil {
		log.Printf("Error saving SMS to database: %s", err)
		http.Error(w, "Failed to save message.", http.StatusInternalServerError)
		return
	}

	// Do something with message
	services.HandleSMS()
}

func handleOutboundSMS(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Failed to parse form data: %s", err)
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	accountSID := os.Getenv("TWILLIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILLIO_AUTH_TOKEN")
	from := os.Getenv("DAVID_TWILIO_PHONE_NUMBER")
	body := r.FormValue("body")
	to := r.FormValue("to")

	twilioURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID)

	formData := url.Values{}
	formData.Set("To", to)
	formData.Set("From", from)
	formData.Set("Body", body)

	userId, err := middleware.GetUserIDFromSession(r)

	if err != nil {
		log.Printf("Error getting user id: %s", err)
		http.Error(w, "Failed to get UserID from session.", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", twilioURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		log.Printf("Error creating request: %s", err)
		http.Error(w, "Failed to create request.", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(accountSID + ":" + authToken))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		http.Error(w, "Failed to send message.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var twilioResp struct {
			SID string `json:"sid"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&twilioResp); err != nil {
			log.Printf("Error parsing Twilio response: %s", err)
			http.Error(w, "Failed to parse Twilio response.", http.StatusInternalServerError)
			return
		}

		sms := models.TextMessage{
			MessageSID: twilioResp.SID,
			UserID:     userId,
			ToNumber:   to,
			FromNumber: from,
			Body:       body,
			Status:     "sent",
			CreatedAt:  time.Now(),
			IsInbound:  false,
		}
		if err := database.SaveSMS(sms); err != nil {
			log.Printf("Error saving SMS to database: %s", err)
			http.Error(w, "Failed to save message to database.", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Error sending request: %s", err)
		http.Error(w, "Failed to send message.", http.StatusInternalServerError)
		return
	}
}
