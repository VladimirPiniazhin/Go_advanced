package main

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/internals/auth"
	"go/order-api/internals/link"
	"go/order-api/internals/product"
	"go/order-api/internals/verify"
	"go/order-api/pkg/db"
	"go/order-api/pkg/middleware"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	db := db.NewDb(config)
	router := http.NewServeMux()

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	productRepository := product.NewProductRepository(db)

	// Handlers
	auth.NewAuthHandler(router)
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
