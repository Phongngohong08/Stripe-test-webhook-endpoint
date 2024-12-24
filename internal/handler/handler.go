package handler

import (
	"encoding/json"
	"fmt"
	"for-learn/stripe-chatgpt/global"
	"for-learn/stripe-chatgpt/internal/models"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/webhook"
)

var productCatalog = map[int]struct {
	Name  string
	Price int64
}{
	1: {"Test Product 1", 2000},
	2: {"Test Product 2", 5000},
	3: {"Test Product 3", 10000},
}

func CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	// Parse product ID from request
	var requestBody struct {
		ProductID int `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get product from catalog
	product, exists := productCatalog[requestBody.ProductID]
	if !exists {
		http.Error(w, "Product not found", http.StatusBadRequest)
		return
	}

	// Create a new Checkout session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(product.Name),
					},
					UnitAmount: stripe.Int64(product.Price),
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

func HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Payment processed successfully")
		// TODO: Xử lý logic thanh toán thành công

	default:
		log.Printf("Unhandled event type: %s", event.Type)
		fmt.Println("Unhandled event type")
	}

	w.WriteHeader(http.StatusOK)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
	}

	// Create a new customer on Stripe
	params := &stripe.CustomerParams{
		Email: stripe.String(req.Username),
	}
	cust, err := customer.New(params)
	if err != nil {
		http.Error(w, "Failed to create Stripe customer", http.StatusInternalServerError)
		return
	}

	user.StripeCustID = cust.ID
	if err := global.Pdb.Save(&user).Error; err != nil {
		http.Error(w, "Failed to save Stripe customer ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "User registered successfully",
		"stripeCustID": cust.ID,
	})
}
