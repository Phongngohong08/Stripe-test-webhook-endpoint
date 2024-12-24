package initialize

import (
	"for-learn/stripe-chatgpt/global"
	"for-learn/stripe-chatgpt/internal/config"
	"for-learn/stripe-chatgpt/internal/models"
	"for-learn/stripe-chatgpt/internal/postgres"
	"for-learn/stripe-chatgpt/internal/router"
	"log"

	"github.com/stripe/stripe-go/v81"
)

func InitAll() {
	config.InitConfig(".env", true)

	// Postgres
	err := postgres.InitPostgres(global.Cfg, &models.Payment{}, &models.Store{}, &models.User{})
	if err != nil {
		log.Fatalf("Error initializing Postgres: %v", err)
	}

	stripe.Key = global.Cfg.STRIPE_SECRET_KEY

	router.InitRouter()
}
