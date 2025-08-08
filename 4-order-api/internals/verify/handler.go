package verify

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/internals/user"
	"go/order-api/pkg/email"
	"go/order-api/pkg/hash"
	"go/order-api/pkg/middleware"
	"go/order-api/pkg/req"
	res "go/order-api/pkg/res"
	"net/http"
)

type VerifyHandlerDeps struct {
	EmailService *email.EmailService
	Config       *configs.Config
}

type VerifyHandler struct {
	EmailService *email.EmailService
	Config       *configs.Config
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &VerifyHandler{
		EmailService: deps.EmailService,
		Config:       deps.Config,
	}

	// Публичный роут (без авторизации)
	router.HandleFunc("GET /verify/{hash}", handler.Verify())

	// Защищённый роут (с авторизацией)
	router.HandleFunc("POST /send", middleware.WithAuth(handler.SendLink(), deps.Config))
}

func (handler *VerifyHandler) SendLink() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[SendLinkRequest](&w, request)
		if err != nil {
			return
		}
		// Проверяем существоввание пользователя
		userData, err := user.FindUserByEmail(body.Address)
		if userData == nil {
			res.JsonResponse(w, 400, "User email is unknown")
			return
		}
		// Генерируем и сохраняем хэш для него если он существует
		verifyHash := hash.GenerateHash(body.Address)
		err = user.SaveHash(userData, verifyHash)
		if err != nil {
			return
		}
		//Отправляем ссылку для верификации для пользователя
		err = handler.EmailService.SendVerificationEmail(body.Address, verifyHash)
		if err != nil {
			return
		}

		result := fmt.Sprintf("Verify link to address %s is sended", body.Address)
		res.JsonResponse(w, 200, result)

	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		// Получаем хеш из URL
		hash := request.PathValue("hash")

		// Получаем пользователя по хешу
		user, err := user.GetUserHash(hash)
		if err != nil {
			res.JsonResponse(w, 500, "Internal server error")
			return
		}
		// Проверяем, найден ли пользователь
		if user == nil {
			res.JsonResponse(w, 400, "Invalid or expired verification link")
			return
		}
		result := fmt.Sprintf("Verification for user email: '%s' is successful", user.Email)
		res.JsonResponse(w, 200, result)
	}
}
