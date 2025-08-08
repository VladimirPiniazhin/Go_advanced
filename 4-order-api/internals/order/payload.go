package order

import (
	"go/order-api/internals/product"
)

type OrderCreateRequest struct {
	UserId   uint              `json:"user_id"`
	Products []product.Product `json:"products"`
}

type OrderUpdateRequest struct {
	UserId   uint              `json:"user_id"`
	Products []product.Product `json:"products"`
}

type OrderGetAllResponse struct {
	Id       uint              `json:"id"`
	UserId   uint              `json:"user_id"`
	Products []product.Product `json:"products"`
}

type OrderGetOneResponse struct {
	Id       uint              `json:"id"`
	UserId   uint              `json:"user_id"`
	Products []product.Product `json:"products"`
}
