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
	err := r.db.
		Preload("CartItems", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC") // หรือ "created_at ASC" ถ้ามี timestamp
		}).
		Preload("CartItems.Product").
		Preload("CartItems.Product.Images").
		Where("user_id = ?", userID).
		FirstOrCreate(&cart, domain.Cart{UserID: userID}).Error

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

func (r *cartRepositoryImpl) DecrementCartItemQuantity(cartID uint, productID uint) error {
	var item domain.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		return err // สินค้าไม่มีในตะกร้า
	}

	if item.Quantity > 1 {
		item.Quantity -= 1
		return r.db.Save(&item).Error
	} else {
		// ถ้าจำนวนเหลือ 1 แล้วลบ => ลบรายการออกจากตะกร้า
		return r.db.Delete(&item).Error
	}
}
