package verify

import (
	"go/order-api/configs"
	"go/order-api/pkg/email"
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
	_ = &VerifyHandler{
		EmailService: deps.EmailService,
		Config:       deps.Config,
	}

}
