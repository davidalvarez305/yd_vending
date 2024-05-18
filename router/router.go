package router

import (
	"net/http"

	"github.com/davidalvarez305/budgeting/handlers"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", handlers.IndexHandler)

	return router
}
