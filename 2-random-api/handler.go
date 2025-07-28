package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

type RandomNumbersHandler struct{}

func NewRandomNumbersHandler(router *http.ServeMux) {
	handler := &RandomNumbersHandler{}
	router.HandleFunc("/random", handler.Random())
}

func (handler *RandomNumbersHandler) Random() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		random := rand.Intn(6) + 1
		fmt.Fprintf(w, "%d", random)
	}
}
