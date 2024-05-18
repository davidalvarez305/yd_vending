package main

import (
	"fmt"
	"net/http"

	"github.com/davidalvarez305/budgeting/router"
)

func main() {
	fmt.Println("Server is listening on port 8000...")
	http.ListenAndServe(":8000", router.Router())
}
