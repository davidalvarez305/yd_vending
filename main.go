package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/davidalvarez305/budgeting/database"
	"github.com/davidalvarez305/budgeting/router"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("ERROR LOADING ENV FILE: %+v\n", err)
	}

	db, err := database.Connect()

	if err != nil {
		log.Fatalf("ERROR CONNECTING TO DB: %+v\n", err)
	}

	fmt.Println("Server is listening on port 8000...")
	http.ListenAndServe(":8000", router.Router())
}
