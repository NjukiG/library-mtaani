package initializers

import (
	"github/NjukiG/library-mtaani/models"
	"log"
)

func SyncDatabase() {
	// Sync the postgrs DB to create tables for the models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Author{},
		&models.Book{},
		&models.Borrow{},
		&models.Cart{},
		&models.CartItem{},
		&models.ShippingDetail{},
		&models.Order{},
		&models.OrderItem{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

}
