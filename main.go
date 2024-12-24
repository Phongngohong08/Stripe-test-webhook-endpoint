package main

import (
	"fmt"
	"for-learn/stripe-chatgpt/global"
	"for-learn/stripe-chatgpt/internal/initialize"
	"log"
	"net/http"
)

func main() {

	initialize.InitAll()

	// Start server
	port := "8080"
	fmt.Printf("Server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, global.Router))
}
