package domain

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`
}

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Quantity    int            `json:"quantity"`
	Images      []ProductImage `gorm:"foreignKey:ProductID" json:"images"`
	CategoryID  uint           `json:"category_id"` // FK
	Category    Category       `json:"category" gorm:"foreignKey:CategoryID"`
}

type ProductImage struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `json:"product_id"`
	Path      string `json:"path"` // local file path เช่น "uploads/1/23/image1.jpg"
}

type ProductRepository interface {
	Create(product Product) error
	GetAllProduct() ([]Product, error)
	GetProductByID(id uint) (*Product, error)
	GetProductByName(name string) (*Product, error)
	UpdateProduct(product Product) error
	Delete(id uint) error
	CreateCategory(category Category) error
	GetProductByNameAndCategoryID(name string, categoryID uint) (*Product, error)

	GetAll() ([]Category, error)

	GetNewArrivals(page, limit int) ([]Product, int64, error)

	GetAllProducts(
		page, limit int,
		sort, order string,
		minPrice, maxPrice float64,
		search string, // เพิ่ม
	) ([]Product, int64, error)

	GetProductByCategory(
		category string,
		page, limit int,
		sort, order string,
		minPrice, maxPrice float64,
	) ([]Product, int64, error)
	CreateBulkProducts(products []*Product) error
}

type ProductService interface {
	// CreateProduct(product Product) error
	GetAllProduct() ([]Product, error)
	UpdateProduct(product Product) error
	GetProductByID(id uint) (*Product, error)
	GetProductByName(name string) (*Product, error)
	Delete(id uint) error
	CreateCategory(category Category) error
	GetAllCategories() ([]Category, error)

	GetNewArrivals(page, limit int) ([]Product, int64, error)
	GetAllProducts(
		page, limit int,
		sort, order string,
		minPrice, maxPrice float64,
		search string, // เพิ่ม
	) ([]Product, int64, error)

	GetProductByCategory(
		category string,
		page, limit int,
		sort, order string,
		minPrice, maxPrice float64,
	) ([]Product, int64, error)
	CreateProducts(products []*Product) ([]*Product, []string, error)
}
