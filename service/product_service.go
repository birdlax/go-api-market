package service

import (
	"backend/domain"
)

type productServiceImpl struct {
	repo domain.ProductRepository
}

func NewProductService(productRepository domain.ProductRepository) domain.ProductService {
	return &productServiceImpl{repo: productRepository}
}

func (s *productServiceImpl) CreateProducts(products []*domain.Product) ([]*domain.Product, []string, error) {
	var toInsert []*domain.Product
	var skipped []string

	for _, p := range products {
		existing, err := s.repo.GetProductByNameAndCategoryID(p.Name, p.CategoryID)
		if err != nil {
			return nil, nil, err
		}
		if existing == nil {
			toInsert = append(toInsert, p)
		} else {
			skipped = append(skipped, p.Name)
		}
	}

	if len(toInsert) > 0 {
		if err := s.repo.CreateBulkProducts(toInsert); err != nil {
			return nil, nil, err
		}
	}

	return toInsert, skipped, nil
}

func (s *productServiceImpl) GetAllProduct() ([]domain.Product, error) {
	products, err := s.repo.GetAllProduct()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productServiceImpl) GetAllProducts(
	page, limit int,
	sort, order string,
	minPrice, maxPrice float64,
	search string, // เพิ่ม
) ([]domain.Product, int64, error) {
	return s.repo.GetAllProducts(page, limit, sort, order, minPrice, maxPrice, search)
}

func (s *productServiceImpl) GetProductByID(id uint) (*domain.Product, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
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

// category service methods
func (s *productServiceImpl) CreateCategory(category domain.Category) error {
	if err := s.repo.CreateCategory(category); err != nil {
		return err
	}
	return nil
}

func (s *productServiceImpl) GetProductByCategory(
	category string,
	page, limit int,
	sort, order string,
	minPrice, maxPrice float64,
) ([]domain.Product, int64, error) {
	return s.repo.GetProductByCategory(category, page, limit, sort, order, minPrice, maxPrice)
}

func (s *productServiceImpl) GetAllCategories() ([]domain.Category, error) {
	return s.repo.GetAll()
}

func (s *productServiceImpl) GetNewArrivals(page, limit int) ([]domain.Product, int64, error) {
	return s.repo.GetNewArrivals(page, limit)
}
