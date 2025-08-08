package order

import (
	"time"
)

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"order_id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
	UserID     uint        `json:"user_id" gorm:"not null"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type OrderItem struct {
	ID        uint       `gorm:"primaryKey" json:"-"` // ID скрыт от JSON
	CreatedAt time.Time  `json:"-"`                   // Скрыт
	UpdatedAt time.Time  `json:"-"`                   // Скрыт
	DeletedAt *time.Time `gorm:"index" json:"-"`      // Скрыт
	OrderID   uint       `json:"-" gorm:"not null"`   // Скрыт (дублируется с order_id)
	ProductID uint       `json:"product_id" gorm:"not null"`
	Quantity  int        `json:"quantity" gorm:"default:1"`
}

func NewOrder(
	userID uint,
	orderItems []OrderItem,
) *Order {

	return &Order{
		UserID:     userID,
		OrderItems: orderItems,
	}
}
