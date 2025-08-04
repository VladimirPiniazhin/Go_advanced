package auth

import (
	"fmt"
	"go/order-api/pkg/req"
	res "go/order-api/pkg/res"
	"net/http"
)

type AuthHandlerDeps struct {
	*AuthService
}

type AuthHandler struct {
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, request)
		if err != nil {
			return
		}
		handler.AuthService.Register(body.Email, body.Password, body.Name)

		result := fmt.Sprintf("Registration user %s with email: %s is successful", body.Name, body.Email)

		res.JsonResponse(w, 201, result)

	}
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		_, err := req.HandleBody[LoginRequest](&w, request)
		if err != nil {
			return
		}
		result := "Login is successful"
		res.JsonResponse(w, 200, result)

	}
}
