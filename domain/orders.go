package domain

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserID        uint        `json:"user_id"`
	TotalPrice    float64     `json:"total_price"`
	Status        string      `json:"status"`
	PaymentMethod string      `json:"payment_method"`
	OrderItems    []OrderItem `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	PaidAt        *time.Time
	AddressID     uint
	Address       Address `gorm:"foreignKey:AddressID"`
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
	GetOrderByID(id uint) (*Order, error)
	UpdateProductStock(tx *gorm.DB, product *Product) error
	UpdateOrder(order Order) (Order, error)
	DeleteOrder(id uint) error
	GetPendingOrderByUserID(userID uint) (Order, error)
	UpdateOrderWithTx(tx *gorm.DB, order *Order) error
	DeleteOrderItemsByOrderID(tx *gorm.DB, orderID uint) error
	CreateOrderItems(tx *gorm.DB, items []OrderItem) error
	GetOrdersByUserIDAndStatus(userID uint, status string) ([]Order, error)
	GetAllOrders(page, limit int, sort, order string) ([]Order, int64, error)
}

type OrderService interface {
	CreateOrder(order Order) (Order, error)
	GetUnpaidOrdersByUserID(userID uint) ([]Order, error)
	GetOrderByID(id uint) (Order, error)
	UpdateOrder(id uint, updated Order) (Order, error)
	DeleteOrder(id uint) error
	MarkOrderAsPaidByUserID(userID uint) error
	CancelOrderByUserID(userID uint) error
	GetOrdersByUserIDAndStatus(userID uint, status string) ([]Order, error)

	GetAllOrders(page, limit int, sort, order string) ([]Order, int64, error)
}
