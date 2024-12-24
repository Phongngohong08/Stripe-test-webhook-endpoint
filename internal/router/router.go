package router

import (
	"for-learn/stripe-chatgpt/global"
	"for-learn/stripe-chatgpt/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRouter() {
	router := mux.NewRouter()

	// Serve static files for the frontend
	router.Handle("/", http.FileServer(http.Dir("./static")))

	// Endpoint to create a Stripe Checkout session
	router.HandleFunc("/create-checkout-session", handler.CreateCheckoutSession).Methods("POST")

	// Endpoint to handle Stripe webhooks
	router.HandleFunc("/webhook", handler.HandleStripeWebhook).Methods("POST")

	// Endpoint to register a new user
	router.HandleFunc("/register", handler.RegisterHandler).Methods("POST")

	global.Router = router
}
