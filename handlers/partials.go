package handlers

import (
	"net/http"

	"github.com/davidalvarez305/yd_vending/constants"
)

func PartialsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/partials/pop-up-modal":
			GetPopUpModal(w, r)
		default:
			http.Error(w, "No partials found.", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetPopUpModal(w http.ResponseWriter, r *http.Request) {
	fileName := "pop_up.html"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PARTIAL_TEMPLATES_DIR+fileName)
}
