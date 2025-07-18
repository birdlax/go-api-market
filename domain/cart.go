package domain

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID    uint       `json:"user_id"`
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	Total     float64    `gorm:"-"`
}

type CartItem struct {
	gorm.Model
	CartID    uint    `json:"cart_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CartItemInput struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type CheckoutRequest struct {
	ShippingAddressID uint   `json:"shipping_address_id"`
	BillingAddressID  uint   `json:"billing_address_id"`
	PaymentMethod     string `json:"payment_method"`
}

// Repository Interface
type CartRepository interface {
	GetCartByUserID(userID uint) (*Cart, error)
	AddOrUpdateCartItem(userID uint, item CartItem) error
	RemoveCartItem(userID, productID uint) error
	ClearCart(userID uint) error
	GetProductByID(productID uint) (Product, error)
	DecrementCartItemQuantity(cartID uint, productID uint) error
}

// Service Interface
type CartService interface {
	AddItem(userID uint, item CartItem) error
	RemoveItem(userID, productID uint) error
	GetCart(userID uint) (*Cart, error)
	GetProductByID(productID uint) (Product, error)
	RemoveItemOne(userID uint, productID uint) error
	AddOneItem(userID uint, productID uint) error
	Checkout(userID uint, req CheckoutRequest) (*Order, error)
}
