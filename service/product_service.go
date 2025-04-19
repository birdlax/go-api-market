package service

import (
	"backend/domain"
)

type ProductService interface {
	Create(product domain.Product) error
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

func (s *productServiceImpl) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
