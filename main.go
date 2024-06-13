package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/middleware"
	"github.com/davidalvarez305/yd_vending/router"
	"github.com/davidalvarez305/yd_vending/sessions"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("ERROR LOADING ENV FILE: %+v\n", err)
	}
	fmt.Println("Environment loaded.")

	_, err = database.Connect()

	if err != nil {
		log.Fatalf("ERROR CONNECTING TO DB: %+v\n", err)
	}
	fmt.Println("Database connected.")

	err = sessions.InitializeSessions()

	if err != nil {
		log.Fatalf("ERROR INITIALIZING SESSIONS: %+v\n", err)
	}
	fmt.Println("Sessions initialized.")

	s := &http.Server{
		Addr:           ":" + constants.ServerPort,
		Handler:        middleware.UserTracking(middleware.SecurityMiddleware(middleware.CSRFProtectMiddleware(router.Router()))),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Server is listening on port 8000...")
	log.Fatal(s.ListenAndServe())
}
