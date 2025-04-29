package repository

import (
	"backend/domain"
	"errors"

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

func (r *productRepositoryImpl) GetAllProduct() ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepositoryImpl) GetProductByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) GetProductByName(name string) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Where("name = ?", name).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) UpdateProduct(product domain.Product) error {
	if err := r.db.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&domain.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}

// category repository
func (r *productRepositoryImpl) CreateCategory(category domain.Category) error {
	if err := r.db.Create(&category).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepositoryImpl) GetProductByCategory(category string) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.
		Preload("Category").
		Where("category_id = ?", category).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepositoryImpl) GetProductByNameAndCategoryID(name string, categoryID uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("name = ? AND category_id = ?", name, categoryID).First(&product).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepositoryImpl) GetAllProducts() ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
