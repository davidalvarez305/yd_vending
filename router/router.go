package router

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/davidalvarez305/yd_vending/handlers"
	"github.com/davidalvarez305/yd_vending/middleware"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()

	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("ERROR GETTING CURRENT DIRECTORY: %+v\n", err)
	}

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(currentDir, "static")))))

	router.Handle("/crm/", middleware.AuthRequired(http.HandlerFunc(handlers.CRMHandler)))
	router.Handle("/inventory/", middleware.AuthRequired(http.HandlerFunc(handlers.InventoryHandler)))
	router.Handle("/external/", middleware.AuthRequired(http.HandlerFunc(handlers.ExternalPagesHandler)))
	router.HandleFunc("/partials/", handlers.PartialsHandler)
	router.HandleFunc("/sms/", handlers.PhoneServiceHandler)
	router.HandleFunc("/call/", handlers.PhoneServiceHandler)
	router.HandleFunc("/webhooks/", handlers.WebhookHandler)
	router.HandleFunc("/funnel/", handlers.FunnelHandler)
	router.HandleFunc("/", handlers.WebsiteHandler)

	return router
}
