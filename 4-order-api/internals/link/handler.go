package link

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/pkg/event"
	"go/order-api/pkg/middleware"
	"go/order-api/pkg/req"
	"go/order-api/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
	Config         *configs.Config
}

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
	Config         *configs.Config
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
		Config:         deps.Config,
	}

	// Публичные роуты (без авторизации)
	router.HandleFunc("GET /{hash}", handler.GoTo())

	// Защищённые роуты (с авторизацией)
	router.HandleFunc("POST /link", middleware.WithAuth(handler.Create(), deps.Config))
	router.HandleFunc("PATCH /link/{id}", middleware.WithAuth(handler.Update(), deps.Config))
	router.HandleFunc("DELETE /link/{id}", middleware.WithAuth(handler.Delete(), deps.Config))
	router.HandleFunc("GET /link", middleware.WithAuth(handler.GetAll(), deps.Config))
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[LinkCreateRequest](&w, request)
		if err != nil {
			return
		}
		newLink := NewLink(body.Url, handler.LinkRepository)
		createdLink, err := handler.LinkRepository.Create(newLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.JsonResponse(w, 201, createdLink)

	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		email, ok := request.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println(email)

		body, err := req.HandleBody[LinkUpdateRequest](&w, request)
		if err != nil {
			return
		}
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result := fmt.Sprintf("Link: %s updated successfully", link.Url)

		res.JsonResponse(w, 201, result)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, 200, "Link deleted successfully")
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		hash := request.PathValue("hash")
		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})
		http.Redirect(w, request, link.Url, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid query parameters", http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(request.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid query parameters", http.StatusBadRequest)
			return
		}
		links := handler.LinkRepository.GetAll(limit, offset)
		count := handler.LinkRepository.Count()
		result := GetAllLinksResponse{
			Links: links,
			Count: count,
		}

		res.JsonResponse(w, 200, result)

	}
}
