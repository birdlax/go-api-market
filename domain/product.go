package domain

import ()

import "time"

type Product struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Quantity    int        `json:"quantity"`
	ImageURL    string     `json:"image_url"`
	Category    string     `json:"category"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
type ProductRepository interface {
	Create(product Product) error
	Delete(id uint) error
}
