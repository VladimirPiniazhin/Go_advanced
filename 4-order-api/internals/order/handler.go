package order

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/pkg/middleware"
	"go/order-api/pkg/req"
	"go/order-api/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type OrderHandlerDeps struct {
	OrderRepository *OrderRepository
	Config          *configs.Config
}

type OrderHandler struct {
	OrderRepository *OrderRepository
	Config          *configs.Config
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderRepository: deps.OrderRepository,
		Config:          deps.Config,
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
		orders, err := handler.OrderRepository.GetAll()
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
		order, err := handler.OrderRepository.GetByID(uint(id))
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
		newOrder := NewOrder(
			body.UserId,
			body.Products,
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
		order, err := handler.OrderRepository.Update(&Order{
			Model:    gorm.Model{ID: uint(id)},
			UserId:   body.UserId,
			Products: body.Products,
		})
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
