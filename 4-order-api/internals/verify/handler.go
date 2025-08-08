package verify

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/internals/user"
	"go/order-api/pkg/email"
	"go/order-api/pkg/hash"
	"go/order-api/pkg/jwt"
	"go/order-api/pkg/req"
	res "go/order-api/pkg/res"
	"net/http"
)

type VerifyHandlerDeps struct {
	EmailService   *email.EmailService
	Config         *configs.Config
	UserRepository *user.UserRepository
}

type VerifyHandler struct {
	EmailService   *email.EmailService
	Config         *configs.Config
	UserRepository *user.UserRepository
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &VerifyHandler{
		EmailService:   deps.EmailService,
		Config:         deps.Config,
		UserRepository: deps.UserRepository,
	}

	// Публичный роут (без авторизации)
	router.HandleFunc("POST /send", handler.SendLink())
	router.HandleFunc("POST /verify/{hash}", handler.Verify())
}

func (handler *VerifyHandler) SendLink() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[SendLinkRequest](&w, request)
		if err != nil {
			return
		}
		// Проверяем существоввание пользователя
		user, err := handler.UserRepository.FindByEmail(body.Address)
		if user == nil {
			res.JsonResponse(w, 400, "User email is unknown")
			return
		}
		// Генерируем и сохраняем хэш для него если он существует
		verifyHash := hash.GenerateHash(body.Address)
		user.Hash = verifyHash
		handler.UserRepository.PatchUser(user)
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
		user, err := handler.UserRepository.FindByHash(hash)
		if err != nil {
			res.JsonResponse(w, 500, "Internal server error")
			return
		}
		// Проверяем, найден ли пользователь
		if user == nil {
			res.JsonResponse(w, 400, "Invalid or expired verification link")
			return
		}
		// Удаляем хеш из пользователя чтобы он не смог его использовать повторно
		user.Hash = ""
		handler.UserRepository.PatchUser(user)
		// Генерируем токен для пользователя
		result := fmt.Sprintf("Verification for user email: '%s' is successful", user.Email)
		token, err := handler.Config.Jwt.Create(jwt.JWTData{
			Phone: user.Phone,
		})
		if err != nil {
			res.JsonResponse(w, 500, "Internal server error")
			return
		}

		res.JsonResponse(w, 200, VerifyResponse{
			Message: result,
			Token:   token,
		})
	}
}
