package user

import (
	"go/order-api/pkg/res"
	"net/http"
	"strconv"
)

type UserHandlerDeps struct {
	UserRepository *UserRepository
}

type UserHandler struct {
	UserRepository *UserRepository
}

func NewUserHandler(router *http.ServeMux, deps UserHandlerDeps) {
	handler := &UserHandler{
		UserRepository: deps.UserRepository,
	}

	router.HandleFunc("DELETE /users/{id}", handler.Delete())

}

func (handler *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.UserRepository.DeleteUser(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, 200, "Product deleted successfully")
	}
}
