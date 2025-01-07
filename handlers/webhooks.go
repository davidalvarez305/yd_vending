package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/types"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		path := r.URL.Path

		if strings.HasPrefix(path, "/webhooks/seed-live-hourly") {
			if len(path) > len("/webhooks/seed-live-hourly") {
				handleSeedLiveHourly(w, r)
				return
			}
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleSeedLiveHourly(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if strings.Contains(strings.ToLower(string(body)), "test") {
		response := map[string]string{
			"status":  "success",
			"message": "OK",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		}
		return
	}

	var transactions []types.SeedLiveTransaction
	if err := json.Unmarshal(body, &transactions); err != nil {
		http.Error(w, "Bad request: Invalid JSON format", http.StatusBadRequest)
		return
	}

	for _, transaction := range transactions {
		err = database.CreateSeedTransaction(transaction)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to save transaction - %s", err), http.StatusInternalServerError)
			return
		}
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Received successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}
