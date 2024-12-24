package config

import (
	"fmt"
	"for-learn/stripe-chatgpt/global"
	"for-learn/stripe-chatgpt/internal/models"
	"log"
	"os"

	"github.com/spf13/viper"
)

func InitConfig(defaultEnvFile string, override bool) {
	cfg := models.Config{}

	env := os.Getenv("ENV")

	if env == "production" {
		fmt.Println("Loading production environment variables")
		cfg.Host = os.Getenv("DB_HOST")
		cfg.Port = os.Getenv("DB_PORT")
		cfg.User = os.Getenv("DB_USER")
		cfg.Password = os.Getenv("DB_PASSWORD")
		cfg.DBName = os.Getenv("DB_NAME")
		cfg.SSLMode = os.Getenv("DB_SSLMODE")
		cfg.STRIPE_SECRET_KEY = os.Getenv("STRIPE_SECRET_KEY")
		cfg.STRIPE_WEBHOOK_SECRET = os.Getenv("STRIPE_WEBHOOK_SECRET")

	} else {
		viper.SetConfigFile(defaultEnvFile)
		if override {
			viper.AutomaticEnv()
		}

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Error reading environment file: %v", err)
		}

		err = viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatalf("Error loading environment file: %v", err)
		}
	}

	global.Cfg = &cfg
}
