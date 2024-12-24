package global

import (
	"for-learn/stripe-chatgpt/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var (
	Pdb    *gorm.DB
	Cfg    *models.Config
	Router *mux.Router
)
