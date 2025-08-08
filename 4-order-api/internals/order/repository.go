package order

import (
	"errors"
	"go/order-api/pkg/db"
)

type OrderRepository struct {
	Database *db.Db
}

func NewOrderRepository(database *db.Db) *OrderRepository {
	return &OrderRepository{
		Database: database,
	}
}

func (repo *OrderRepository) Create(order *Order) (*Order, error) {
	result := repo.Database.DB.Create(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (repo *OrderRepository) GetByID(id uint, userID uint) (*Order, error) {
	var order Order
	result := repo.Database.DB.Preload("OrderItems").First(&order, "id = ? AND user_id = ?", id, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}
func (repo *OrderRepository) GetAll(userID uint) (*[]Order, error) {
	var orders []Order
	result := repo.Database.DB.Preload("OrderItems").Find(&orders, "user_id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &orders, nil
}

func (repo *OrderRepository) Update(order *Order) (*Order, error) {
	result := repo.Database.DB.Updates(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}
func (repo *OrderRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Order{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Url not exist")
	}
	return nil
}
