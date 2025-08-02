package product

import "github.com/lib/pq"

type ProductCreateRequest struct {
	Description string         `json:"description"`
	Name        string         `json:"name"`
	Images      pq.StringArray `json:"img" gorm:"type:text[]"`
	Price       int            `json:"price"`
}

type ProductUpdateRequest struct {
	Description string         `json:"description"`
	Name        string         `json:"name"`
	Images      pq.StringArray `json:"img" gorm:"type:text[]"`
	Price       int            `json:"price"`
}
