package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/models"
	"github.com/davidalvarez305/yd_vending/services"
	"github.com/davidalvarez305/yd_vending/types"
)

func PhoneServiceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/call/inbound":
			handleInboundCall(w, r)
		case "/call/inbound/end":
			handleInboundCallEnd(w, r)
		case "/sms/inbound":
			handleInboundSMS(w, r)
		case "/sms/outbound":
			handleOutboundSMS(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleInboundCall(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	var incomingPhoneCall types.TwilioIncomingCallBody
	if err := decoder.Decode(&incomingPhoneCall, r.PostForm); err != nil {
		fmt.Printf("Error decoding form data: %+v\n", err)
		http.Error(w, "Failed to decode form data", http.StatusBadRequest)
		return
	}

	forwardNumber, err := database.GetForwardPhoneNumber(helpers.RemoveCountryCode(incomingPhoneCall.To), helpers.RemoveCountryCode(incomingPhoneCall.From))
	if err != nil {
		http.Error(w, "Failed to get matching phone number.", http.StatusInternalServerError)
		return
	}

	twiML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Dial action="%s">%s</Dial>
	</Response>`, forwardNumber.ForwardPhoneNumber, constants.TwilioCallbackWebhook)

	phoneCall := models.PhoneCall{
		ExternalID:   incomingPhoneCall.CallSid,
		UserID:       forwardNumber.UserID,
		LeadID:       forwardNumber.LeadID,
		CallDuration: 0,
		DateCreated:  time.Now().Unix(),
		CallFrom:     incomingPhoneCall.From,
		CallTo:       incomingPhoneCall.To,
		IsInbound:    incomingPhoneCall.Direction == "inbound",
		RecordingURL: "",
		Status:       incomingPhoneCall.CallStatus,
	}

	if err := database.SavePhoneCall(phoneCall); err != nil {
		http.Error(w, "Failed to save phone call.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(twiML))
}

func handleInboundCallEnd(w http.ResponseWriter, r *http.Request) {
	var dialStatus types.IncomingPhoneCallDialStatus

	if err := json.NewDecoder(r.Body).Decode(&dialStatus); err != nil {
		http.Error(w, "Failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	phoneCall, err := database.GetPhoneCallBySID(dialStatus.DialCallSid)
	if err != nil {
		http.Error(w, "Failed to get phone call by SID.", http.StatusInternalServerError)
		return
	}

	phoneCall.CallDuration = dialStatus.DialCallDuration
	phoneCall.RecordingURL = dialStatus.RecordingURL

	if err := database.UpdatePhoneCall(phoneCall); err != nil {
		http.Error(w, "Failed to save phone call.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
}

func handleInboundSMS(w http.ResponseWriter, r *http.Request) {
	var twilioMessage types.TwilioMessage

	if err := json.NewDecoder(r.Body).Decode(&twilioMessage); err != nil {
		http.Error(w, "Failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	userId, err := database.GetUserIDFromPhoneNumber(helpers.RemoveCountryCode(twilioMessage.To))
	if err != nil {
		http.Error(w, "Failed to get User ID.", http.StatusBadRequest)
		return
	}

	leadId, err := database.GetLeadIDFromPhoneNumber(helpers.RemoveCountryCode(twilioMessage.From))
	if err != nil {
		http.Error(w, "Failed to get Lead ID.", http.StatusBadRequest)
		return
	}

	dateCreated := time.Unix(twilioMessage.DateCreated.Unix(), 0).Unix()

	message := models.Message{
		ExternalID:  twilioMessage.MessageSid,
		UserID:      userId,
		LeadID:      leadId,
		Text:        twilioMessage.Body,
		TextFrom:    helpers.RemoveCountryCode(twilioMessage.From),
		TextTo:      helpers.RemoveCountryCode(twilioMessage.To),
		IsInbound:   true,
		DateCreated: dateCreated,
		Status:      twilioMessage.SmsStatus,
	}

	if err := database.SaveSMS(message); err != nil {
		log.Printf("Error saving SMS to database: %s", err)
		http.Error(w, "Failed to save message.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleOutboundSMS(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid request.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	form := types.OutboundMessageForm{
		To:   r.FormValue("to"),
		Body: r.FormValue("body"),
		From: r.FormValue("from"),
	}

	userId, err := database.GetUserIDFromPhoneNumber(form.From)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Could not find matching user.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadId, err := database.GetLeadIDFromPhoneNumber(form.To)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Could not find matching lead.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	messageResponse, err := services.SendOutboundMessage(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to send text message.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	message := models.Message{
		ExternalID:  messageResponse.Sid,
		UserID:      userId,
		LeadID:      leadId,
		Text:        form.Body,
		TextFrom:    form.From,
		TextTo:      form.To,
		IsInbound:   false,
		DateCreated: time.Now().Unix(),
		Status:      messageResponse.Status,
	}

	err = database.SaveSMS(message)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to save message.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	messages, err := database.GetMessagesByLeadID(leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get new messages.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "messages.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "messages.html",
		Data: map[string]any{
			"Messages": messages,
		},
	}

	w.WriteHeader(http.StatusOK)
	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
