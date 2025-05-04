package repository

import (
	"backend/domain"
	"gorm.io/gorm"
)

type cartRepositoryImpl struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) domain.CartRepository {
	return &cartRepositoryImpl{db: db}
}

func (r *cartRepositoryImpl) GetCartByUserID(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Preload("CartItems.Product").Where("user_id = ?", userID).FirstOrCreate(&cart, domain.Cart{UserID: userID}).Error
	return &cart, err
}

func (r *cartRepositoryImpl) AddOrUpdateCartItem(userID uint, item domain.CartItem) error {
	cart, err := r.GetCartByUserID(userID)
	if err != nil {
		return err
	}
	item.CartID = cart.ID

	// Check if the item already exists in the cart
	var existingItem domain.CartItem
	err = r.db.Where("cart_id = ? AND product_id = ?", cart.ID, item.ProductID).First(&existingItem).Error

	if err == nil {
		// Update the existing item
		existingItem.Quantity += item.Quantity
		return r.db.Save(&existingItem).Error
	}

	// Add new item to the cart
	return r.db.Create(&item).Error
}

func (r *cartRepositoryImpl) RemoveCartItem(userID, productID uint) error {
	cart, err := r.GetCartByUserID(userID)
	if err != nil {
		return err
	}
	return r.db.Where("cart_id = ? AND product_id = ?", cart.ID, productID).Delete(&domain.CartItem{}).Error
}

func (r *cartRepositoryImpl) ClearCart(userID uint) error {
	cart, err := r.GetCartByUserID(userID)
	if err != nil {
		return err
	}
	return r.db.Where("cart_id = ?", cart.ID).Delete(&domain.CartItem{}).Error
}

func (r *cartRepositoryImpl) GetProductByID(productID uint) (domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, productID).Error
	return product, err
}
