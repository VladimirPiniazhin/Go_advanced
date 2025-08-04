package main

import (
	"go/order-api/internals/link"
	"go/order-api/internals/product"
	"go/order-api/internals/user"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&link.Link{})
	db.AutoMigrate(&product.Product{})
	db.AutoMigrate(&user.User{})
}
