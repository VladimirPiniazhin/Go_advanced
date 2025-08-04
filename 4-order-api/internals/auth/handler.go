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
	router.HandleFunc("POST /auth/sendsms", handler.VerifyBySMS())
	router.HandleFunc("POST /auth/verify", handler.GetTokenByPhoneVerification())
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, request)
		if err != nil {
			return
		}
		token, err := handler.AuthService.Register(body.Email, body.Password, body.Name, body.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		result := fmt.Sprintf("Registration user %s with email: %s is successful. Token: %s", body.Name, body.Email, token)

		res.JsonResponse(w, 201, result)

	}
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		body, err := req.HandleBody[LoginRequest](&w, request)
		if err != nil {

			return
		}

		token, err := handler.AuthService.UserLogin(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		result := AuthorizationResponse{
			Token: token,
			Msg:   fmt.Sprintf("User: %s -  login is successful.", body.Email),
		}
		res.JsonResponse(w, 200, result)

	}
}

func (handler *AuthHandler) VerifyBySMS() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		body, err := req.HandleBody[VerifyRequest](&w, request)
		if err != nil {
			return
		}

		sessionID, code, err := handler.AuthService.UpdateSessionID(body.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := VerifyResponse{
			Session: Session{
				SessionID: sessionID,
				Code:      code,
			},
		}
		res.JsonResponse(w, 200, result)

	}
}

func (handler *AuthHandler) GetTokenByPhoneVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		body, err := req.HandleBody[AuthorizationBySMSRequest](&w, request)
		if err != nil {
			return
		}

		token, err := handler.VerifyUserBySmsCode(body.SessionID, body.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := AuthorizationResponse{
			Token: token,
			Msg:   "Verification by sms is successful",
		}

		res.JsonResponse(w, 200, result)

	}
}
