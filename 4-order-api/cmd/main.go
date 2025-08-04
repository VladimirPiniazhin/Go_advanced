package main

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/internals/auth"
	"go/order-api/internals/link"
	"go/order-api/internals/product"
	"go/order-api/internals/user"
	"go/order-api/internals/verify"
	"go/order-api/pkg/db"
	"go/order-api/pkg/middleware"
	"log"
	"net/http"
	"os"

	logger "github.com/sirupsen/logrus"
)

func main() {
	config := configs.LoadConfig()
	db := db.NewDb(config)
	router := http.NewServeMux()

	// Logging
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов: %v", err)
	}
	logger.SetOutput(file)
	logger.SetFormatter(&logger.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetLevel(logger.InfoLevel)

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	productRepository := product.NewProductRepository(db)
	userRepository := user.NewUserRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})
	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		Config: config,
	})
	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})

	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
		middleware.IsAuthed,
	)
	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}
	fmt.Println("Server listening on port 8081")
	server.ListenAndServe()

}
