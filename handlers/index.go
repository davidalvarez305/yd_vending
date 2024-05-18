package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/davidalvarez305/budgeting/constants"
	"github.com/davidalvarez305/budgeting/helpers"
	"github.com/davidalvarez305/budgeting/models"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expense := models.Expense{
		ExpenseAmount:   parseInt(r.Form.Get("expenseAmount")),
		ExpenseCategory: parseInt(r.Form.Get("expenseCategory")),
		ExpenseDate:     parseDateTime(r.Form.Get("expenseDate")),
		ExpenseComments: r.Form.Get("expenseComments"),
	}

	fileName := "form.html"

	err = helpers.BuildFile(fileName, constants.PUBLIC_DIR+fileName, constants.TEMPLATES_DIR+fileName, expense)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PUBLIC_DIR+fileName)
}

func parseInt(s string) int {
	i := 0
	if s != "" {
		i, _ = strconv.Atoi(s)
	}
	return i
}

func parseDateTime(s string) time.Time {
	layout := "2006-01-02T15:04" // datetime-local input format
	t, _ := time.Parse(layout, s)
	return t
}
