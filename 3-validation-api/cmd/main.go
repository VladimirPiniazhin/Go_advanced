package main

import (
	"fmt"
	config "go/verify-api/configs"
	"go/verify-api/internals/verify"
	"net/http"
)

func main() {
	config := config.LoadConfig()
	router := http.NewServeMux()
	verify.NewVerifyHandler(router, &verify.VerifyHandlerDeps{
		Config: config,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	fmt.Println("Server listening on port 8081")
	server.ListenAndServe()

}
