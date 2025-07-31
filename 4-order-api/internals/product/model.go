package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Description string         `json:"description"`
	Name        string         `json:"name"`
	Image       pq.StringArray `json:"img" gorm:"type:text[]"`
}

type Cart struct {
	Cart []Product `json:"cart"`
}
