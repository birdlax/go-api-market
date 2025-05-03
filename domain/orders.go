package domain

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	PaidAt     *time.Time
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderRepository interface {
	CreateWithTx(tx *gorm.DB, order *Order) error
	BeginTx() *gorm.DB

	GetProductByID(id uint) (*Product, error)
	GetAllOrders() ([]Order, error)
	GetOrderByID(id uint) (*Order, error)
	UpdateProductStock(tx *gorm.DB, product *Product) error
	UpdateOrder(order Order) (Order, error)
	DeleteOrder(id uint) error
}
