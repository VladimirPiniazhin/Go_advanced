package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Description string         `json:"description"`
	Name        string         `json:"name"`
	Images      pq.StringArray `json:"img" gorm:"type:text[]"`
	Price       int            `json:"price"`
}

func NewProduct(
	description string,
	name string,
	images pq.StringArray,
	price int,
	linkRepo *ProductRepository) *Product {

	return &Product{
		Description: description,
		Name:        name,
		Images:      images,
		Price:       price,
	}
}
