package service

import (
	"backend/domain"
	"fmt"
)

type cartServiceImpl struct {
	repo         domain.CartRepository
	orderService domain.OrderService
}

func NewCartService(repo domain.CartRepository, orderService domain.OrderService) domain.CartService {
	return &cartServiceImpl{repo: repo, orderService: orderService}
}

func (s *cartServiceImpl) AddItem(userID uint, item domain.CartItem) error {
	product, err := s.repo.GetProductByID(item.ProductID)
	if err != nil {
		return err
	}
	item.Price = product.Price // คำนวณราคาสินค้า (ถ้าจำเป็น)

	// ตรวจสอบว่ามีสินค้านี้ในตะกร้าหรือไม่
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	for _, cartItem := range cart.CartItems {
		if cartItem.ProductID == item.ProductID {
			item.Quantity += cartItem.Quantity // เพิ่มจำนวนสินค้าในตะกร้า
			break
		}
	}

	return s.repo.AddOrUpdateCartItem(userID, item)
}

func (s *cartServiceImpl) RemoveItem(userID, productID uint) error {
	return s.repo.RemoveCartItem(userID, productID)
}

func (s *cartServiceImpl) GetCart(userID uint) (*domain.Cart, error) {
	return s.repo.GetCartByUserID(userID)
}

func (s *cartServiceImpl) Checkout(userID uint) error {
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return err
	}
	if len(cart.CartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	// เตรียม Order
	order := domain.Order{
		UserID:     userID,
		OrderItems: make([]domain.OrderItem, len(cart.CartItems)),
	}

	for i, item := range cart.CartItems {
		order.OrderItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	_, err = s.orderService.CreateOrder(order)
	if err != nil {
		return err
	}

	// ล้างตะกร้า
	if err := s.repo.ClearCart(userID); err != nil {
		return err
	}

	return nil
}

func (s *cartServiceImpl) GetProductByID(productID uint) (domain.Product, error) {
	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}
