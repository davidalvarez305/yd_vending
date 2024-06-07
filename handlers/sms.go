package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/middleware"
	"github.com/davidalvarez305/yd_vending/models"
	"github.com/davidalvarez305/yd_vending/services"
	"github.com/davidalvarez305/yd_vending/types"
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
	fmt.Println(r.FormValue("Body"))
}

func handleOutboundSMS(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Failed to parse form data: %s", err)
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	form := types.OutboundMessageForm{
		To:   r.FormValue("to"),
		Body: r.FormValue("Body"),
		From: os.Getenv("DAVID_TWILIO_PHONE_NUMBER"),
	}

	userId, err := middleware.GetUserIDFromSession(r)
	if err != nil {
		log.Printf("Error getting user id: %s", err)
		http.Error(w, "Failed to get UserID from session.", http.StatusInternalServerError)
		return
	}

	messageSID, err := services.SendOutboundMessage(form)
	if err != nil {
		log.Printf("Error sending message: %s", err)
		http.Error(w, "Error sending message.", http.StatusInternalServerError)
		return
	}

	sms := models.TextMessage{
		MessageSID: messageSID,
		UserID:     userId,
		ToNumber:   form.To,
		FromNumber: form.From,
		Body:       form.Body,
		Status:     "sent",
		CreatedAt:  time.Now(),
		IsInbound:  false,
	}
	if err := database.SaveSMS(sms); err != nil {
		log.Printf("Error saving SMS to database: %s", err)
		http.Error(w, "Failed to save message to database.", http.StatusInternalServerError)
		return
	}
}
