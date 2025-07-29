package auth

import (
	"fmt"
	"go/http_serv/configs"
	"go/http_serv/pkg/req"
	res "go/http_serv/pkg/res"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
}

type AuthHandler struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps *AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
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
		result := fmt.Sprintf("Registration user %s is successful", body.Name)

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
