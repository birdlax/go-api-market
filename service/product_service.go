package service

import (
	"backend/domain"
	"errors"
)

type ProductService interface {
	CreateProduct(product domain.Product) error
	CreateCategory(category domain.Category) error
	GetAllProduct() ([]domain.Product, error)
	UpdateProduct(product domain.Product) error
	GetProductByName(name string) (*domain.Product, error)
	GetProductByCategory(category string) (*domain.Product, error)
	Delete(id uint) error
}

type productServiceImpl struct {
	repo domain.ProductRepository
}

func NewProductService(productRepository domain.ProductRepository) ProductService {
	return &productServiceImpl{repo: productRepository}
}

func (s *productServiceImpl) CreateProduct(product domain.Product) error {
	existingProduct, err := s.repo.GetProductByNameAndCategoryID(product.Name, product.CategoryID)
	if err != nil {
		return err
	}
	if existingProduct != nil {
		return errors.New("product already exists in this category")
	}

	return s.repo.Create(product)
}

func (s *productServiceImpl) CreateCategory(category domain.Category) error {
	if err := s.repo.CreateCategory(category); err != nil {
		return err
	}
	return nil
}

func (s *productServiceImpl) GetAllProduct() ([]domain.Product, error) {
	products, err := s.repo.GetAllProduct()
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (s *productServiceImpl) GetProductByName(name string) (*domain.Product, error) {
	product, err := s.repo.GetProductByName(name)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productServiceImpl) GetProductByCategory(category string) (*domain.Product, error) {
	product, err := s.repo.GetProductByCategory(category)
	if err != nil {
		return nil, err
	}
	return product, nil
}
func (s *productServiceImpl) UpdateProduct(product domain.Product) error {
	if err := s.repo.UpdateProduct(product); err != nil {
		return err
	}
	return nil

}

func (s *productServiceImpl) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
