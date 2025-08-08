package order

import (
	"go/order-api/internals/product"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserId   uint              `json:"user_id"`
	Products []product.Product `json:"products"`
}

func NewOrder(
	userId uint,
	products []product.Product,
) *Order {

	return &Order{
		UserId:   userId,
		Products: products,
	}
}
