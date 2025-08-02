package main

import (
	"fmt"
	"go/http_serv/configs"
	"go/http_serv/internals/auth"
	"go/http_serv/internals/randomnums"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	randomnums.NewRandomNumbersHandler(router)
	auth.NewAuthHandler(router, &auth.AuthHandlerDeps{
		Config: config,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	fmt.Println("Server listening on port 8081")
	server.ListenAndServe()

}
