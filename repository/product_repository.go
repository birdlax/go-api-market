package repository

import (
	"backend/domain"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(product domain.Product) error {
	return r.db.Create(&product).Error
}

func (r *productRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&domain.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}
