package order

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/internals/user"
	"go/order-api/pkg/middleware"
	"go/order-api/pkg/req"
	"go/order-api/pkg/res"
	"net/http"
	"strconv"
)

type OrderHandlerDeps struct {
	OrderRepository *OrderRepository
	Config          *configs.Config
	UserRepository  *user.UserRepository
}

type OrderHandler struct {
	OrderRepository *OrderRepository
	Config          *configs.Config
	UserRepository  *user.UserRepository
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderRepository: deps.OrderRepository,
		Config:          deps.Config,
		UserRepository:  deps.UserRepository,
	}

	// Защищённые роуты (с авторизацией)
	router.HandleFunc("GET /my-orders", middleware.WithAuth(handler.GetAll(), deps.Config))
	router.HandleFunc("GET /order/{id}", middleware.WithAuth(handler.GetOne(), deps.Config))
	router.HandleFunc("POST /order", middleware.WithAuth(handler.Create(), deps.Config))
	router.HandleFunc("PATCH /order/{id}", middleware.WithAuth(handler.Update(), deps.Config))
	router.HandleFunc("DELETE /order/{id}", middleware.WithAuth(handler.Delete(), deps.Config))
}

func (handler *OrderHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		phone, ok := request.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := handler.UserRepository.FindByPhoneNumber(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		orders, err := handler.OrderRepository.GetAll(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		res.JsonResponse(w, http.StatusOK, orders)
	}
}

func (handler *OrderHandler) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		phone, ok := request.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := handler.UserRepository.FindByPhoneNumber(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		order, err := handler.OrderRepository.GetByID(uint(id), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		res.JsonResponse(w, http.StatusOK, order)
	}
}

func (handler *OrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[OrderCreateRequest](&w, request)
		if err != nil {
			return
		}
		phone, ok := request.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := handler.UserRepository.FindByPhoneNumber(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		newOrder := NewOrder(
			user.ID,
			body.OrderItems,
		)
		createdOrder, err := handler.OrderRepository.Create(newOrder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.JsonResponse(w, 201, createdOrder)

	}
}

func (handler *OrderHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[OrderUpdateRequest](&w, request)
		if err != nil {
			return
		}
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		phone, ok := request.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := handler.UserRepository.FindByPhoneNumber(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		order, err := handler.OrderRepository.GetByID(uint(id), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		order.OrderItems = body.OrderItems
		order, err = handler.OrderRepository.Update(order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result := fmt.Sprintf("Order: %d updated successfully", order.ID)

		res.JsonResponse(w, 201, result)
	}
}

func (handler *OrderHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.OrderRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, 200, "Order deleted successfully")
	}
}
