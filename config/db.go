package config

import (
	"backend/domain"
	"fmt"
	// "github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("❌ Error loading .env file")
	// }

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("❌ Missing one or more environment variables for database connection.")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Cart{},
		&domain.CartItem{},
		&domain.Address{},
	)
	if err != nil {
		log.Fatal("❌ Failed to migrate:", err)
	}

	DB = db
	fmt.Println("✅ Database connected and migrated successfully.")
}
