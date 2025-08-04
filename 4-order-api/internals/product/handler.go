package product

import (
	"fmt"
	"go/order-api/pkg/req"
	"go/order-api/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	ProductRepository *ProductRepository
}

type ProductHandler struct {
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
	}
	router.HandleFunc("GET /products", handler.GetAll())
	router.HandleFunc("GET /products/{id}", handler.GetOne())
	router.HandleFunc("POST /products", handler.Create())
	router.HandleFunc("PATCH /products/{id}", handler.Update())
	router.HandleFunc("DELETE /products/{id}", handler.Delete())

}

func (handler *ProductHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		products, err := handler.ProductRepository.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		res.JsonResponse(w, http.StatusOK, products)
	}
}

func (handler *ProductHandler) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product, err := handler.ProductRepository.GetByID(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		res.JsonResponse(w, http.StatusOK, product)
	}
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[ProductCreateRequest](&w, request)
		if err != nil {
			return
		}
		newProduct := NewProduct(
			body.Description,
			body.Name,
			body.Images,
			body.Price,
			handler.ProductRepository)
		createdProduct, err := handler.ProductRepository.Create(newProduct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.JsonResponse(w, 201, createdProduct)

	}
}

func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[ProductUpdateRequest](&w, request)
		if err != nil {
			return
		}
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product, err := handler.ProductRepository.Update(&Product{
			Model:       gorm.Model{ID: uint(id)},
			Description: body.Description,
			Name:        body.Name,
			Images:      body.Images,
			Price:       body.Price,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result := fmt.Sprintf("Product: %s updated successfully", product.Name)

		res.JsonResponse(w, 201, result)
	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.ProductRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, 200, "Product deleted successfully")
	}
}
