package helpers

import (
	"errors"
	"html/template"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/davidalvarez305/budgeting/models"
)

func BuildFile(fileName, publicFilePath, templateFilePath string, data any) error {
	tmpl, err := template.New(fileName).ParseFiles(templateFilePath)

	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Create(publicFilePath)

	if err != nil {
		return err
	}

	err = tmpl.Execute(f, data)

	if err != nil {
		return err
	}

	err = f.Close()

	if err != nil {
		return err
	}

	return err
}

func parseInt(s string) (int, error) {
	i := 0
	if s != "" {
		return strconv.Atoi(s)
	}
	return i, errors.New("int cannot be empty string")
}

func parseDateTime(s string) (time.Time, error) {
	layout := "2006-01-02T15:04" // datetime-local input format
	return time.Parse(layout, s)
}

func parseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func ParseTransaction(form url.Values) (models.Transaction, error) {
	var transaction models.Transaction

	amount, err := parseInt(form.Get("amount"))

	if err != nil {
		return transaction, err
	}

	category, err := parseInt(form.Get("category"))

	if err != nil {
		return transaction, err
	}

	createdAt, err := parseDateTime(form.Get("category"))

	if err != nil {
		return transaction, err
	}

	isFixed, err := parseBool(form.Get("is_fixed"))

	if err != nil {
		return transaction, err
	}

	isExpense, err := parseBool(form.Get("is_expense"))

	if err != nil {
		return transaction, err
	}

	transaction = models.Transaction{
		Amount:     amount,
		CategoryID: category,
		CreatedAt:  createdAt,
		Comments:   form.Get("comments"),
		IsFixed:    isFixed,
		IsExpense:  isExpense,
	}

	return transaction, err
}
