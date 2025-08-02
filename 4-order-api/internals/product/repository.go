package product

import (
	"errors"
	"go/order-api/pkg/db"
)

type ProductRepository struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}

func (repo *ProductRepository) Create(product *Product) (*Product, error) {
	result := repo.Database.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *ProductRepository) GetByID(id uint) (*Product, error) {
	var product Product
	result := repo.Database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}
func (repo *ProductRepository) GetAll() (*[]Product, error) {
	var products []Product
	result := repo.Database.DB.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return &products, nil
}

func (repo *ProductRepository) Update(product *Product) (*Product, error) {
	result := repo.Database.DB.Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}
func (repo *ProductRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Url not exist")
	}
	return nil
}
