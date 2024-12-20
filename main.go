package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"
)

func main() {
	// Load environment variables
	err := godotenv.Load(os.Getenv("ENV_FILE_PATH"))
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set Stripe API Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Create a new router
	router := mux.NewRouter()

	// Serve static files for the frontend
	router.Handle("/", http.FileServer(http.Dir("./static")))

	// Endpoint to create a Stripe Checkout session
	router.HandleFunc("/create-checkout-session", createCheckoutSession).Methods("POST")

	// Endpoint to handle Stripe webhooks
	router.HandleFunc("/webhook", handleStripeWebhook).Methods("POST")

	// Start server
	port := "8080"
	fmt.Printf("Server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func createCheckoutSession(w http.ResponseWriter, r *http.Request) {
	// Create a new Checkout session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("This is test Product"),
					},
					UnitAmount: stripe.Int64(int64(2000)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("http://localhost:8080/success.html"),
		CancelURL:  stripe.String("http://localhost:8080/cancel.html"),
	}

	sess, err := session.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Checkout Session: %v", err), http.StatusInternalServerError)
		return
	}

	// Return session ID to the frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id": sess.ID,
	})
}

func handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusServiceUnavailable)
		return
	}

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event := stripe.Event{}

	sigHeader := r.Header.Get("Stripe-Signature")

	event, err = webhook.ConstructEvent(payload, sigHeader, webhookSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Webhook signature verification failed: %v", err), http.StatusBadRequest)
		fmt.Println("error: ", err)
		fmt.Println("Webhook signature verification failed")
		return
	}

	// Xử lý các event Stripe gửi về
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			http.Error(w, "Error parsing webhook JSON", http.StatusBadRequest)
			return
		}

		log.Printf("Payment succeeded: %s", session.ID)
		fmt.Println("Payment processed successfully") // Chỉ in ra khi không gặp lỗi
		// TODO: Xử lý logic thanh toán thành công

	default:
		log.Printf("Unhandled event type: %s", event.Type)
		fmt.Println("Unhandled event type")
	}

	// Đảm bảo phản hồi thành công được ghi một lần
	w.WriteHeader(http.StatusOK)
}
