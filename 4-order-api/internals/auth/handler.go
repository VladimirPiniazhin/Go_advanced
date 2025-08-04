package auth

import (
	"fmt"
	"go/order-api/internals/user"
	"go/order-api/pkg/req"
	res "go/order-api/pkg/res"
	"net/http"
)

type AuthHandlerDeps struct {
}

type AuthHandler struct {
}

func NewAuthHandler(router *http.ServeMux) {
	handler := &AuthHandler{}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, request)
		if err != nil {
			return
		}
		err = user.SaveUser(body.Email, body.Password, body.Name)
		if err != nil {
			res.JsonResponse(w, http.StatusBadRequest, err.Error())
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
