package verify

import (
	config "go/verify-api/configs"
	res "go/verify-api/pkg/res"
	"net/http"
)

type VerifyHandlerDeps struct {
	*config.Config
}

type VerifyHandler struct {
	*config.Config
}

func NewVerifyHandler(router *http.ServeMux, deps *VerifyHandlerDeps) {
	handler := &VerifyHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /send", handler.SendLink())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
}

func (handler *VerifyHandler) SendLink() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := VerifyResponse{
			Msg: "Link is sended",
		}
		res.JsonRes(w, 200, data)
	}
}
func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}
