package service

import (
	"backend/domain"
)

type ProductService interface {
	Create(product domain.Product) error
	GetAllProduct() ([]domain.Product, error)
	UpdateProduct(product domain.Product) error
	GetProductByName(name string) (*domain.Product, error)
	Delete(id uint) error
}

type productServiceImpl struct {
	repo domain.ProductRepository
}

func NewProductService(productRepository domain.ProductRepository) ProductService {
	return &productServiceImpl{repo: productRepository}
}

func (s *productServiceImpl) Create(product domain.Product) error {
	if err := s.repo.Create(product); err != nil {
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
