package handlers

import (
	"net/http"

	"github.com/davidalvarez305/budgeting/constants"
	"github.com/davidalvarez305/budgeting/helpers"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetIndex(w, r)
	case http.MethodPost:
		PostIndex(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	fileName := "index.html"

	err := helpers.BuildFile(fileName, constants.PUBLIC_DIR+fileName, constants.TEMPLATES_DIR+fileName, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PUBLIC_DIR+fileName)
}

func PostIndex(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Form cannot be parsed.", http.StatusBadRequest)
		return
	}

	transaction, err := helpers.ParseTransaction(r.Form)

	if err != nil {
		http.Error(w, "Error parsing transaction.", http.StatusInternalServerError)
		return
	}

	fileName := "form.html"

	err = helpers.BuildFile(fileName, constants.PUBLIC_DIR+fileName, constants.TEMPLATES_DIR+fileName, transaction)

	if err != nil {
		http.Error(w, "Error building HTML file.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PUBLIC_DIR+fileName)
}
