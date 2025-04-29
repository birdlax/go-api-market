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
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Quantity    int      `json:"quantity"`
	ImageURL    string   `json:"image_url"`
	CategoryID  uint     `json:"category_id"` // FK
	Category    Category `json:"category" gorm:"foreignKey:CategoryID"`
}

type ProductRepository interface {
	Create(product Product) error
	GetAllProduct() ([]Product, error)
	GetProductByID(id uint) (*Product, error)
	GetProductByName(name string) (*Product, error)
	UpdateProduct(product Product) error
	Delete(id uint) error

	GetAllProducts() ([]Product, error)
	CreateCategory(category Category) error
	GetProductByCategory(category string) ([]Product, error)
	GetProductByNameAndCategoryID(name string, categoryID uint) (*Product, error)
}
