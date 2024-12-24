package postgres

import (
	"fmt"
	"for-learn/stripe-chatgpt/global"
	"for-learn/stripe-chatgpt/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPostgres(cfg *models.Config, models ...interface{}) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Postgres connected")

	// Perform auto-migration for provided models
	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			return fmt.Errorf("failed to auto-migrate models: %w", err)
		}
	}

	// Assign the initialized database to the global variable
	global.Pdb = db
	return nil
}
