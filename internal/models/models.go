package models

type Payment struct {
	ID          string `json:"id"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

type Store struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	Password     string `gorm:"not null"`
	StripeCustID string `gorm:"unique"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Host                  string `mapstructure:"DB_HOST"`
	Port                  string `mapstructure:"DB_PORT"`
	User                  string `mapstructure:"DB_USER"`
	Password              string `mapstructure:"DB_PASSWORD"`
	DBName                string `mapstructure:"DB_NAME"`
	SSLMode               string `mapstructure:"DB_SSLMODE"`
	STRIPE_SECRET_KEY     string `mapstructure:"STRIPE_SECRET_KEY"`
	STRIPE_WEBHOOK_SECRET string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
}
