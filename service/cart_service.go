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

	return s.repo.AddOrUpdateCartItem(userID, item)
}

func (s *cartServiceImpl) RemoveItem(userID, productID uint) error {
	return s.repo.RemoveCartItem(userID, productID)
}

func (s *cartServiceImpl) GetCart(userID uint) (*domain.Cart, error) {
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	var total float64
	for _, item := range cart.CartItems {
		total += float64(item.Quantity) * item.Product.Price
	}
	cart.Total = total

	return cart, nil
}

func (s *cartServiceImpl) GetProductByID(productID uint) (domain.Product, error) {
	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (s *cartServiceImpl) RemoveItemOne(userID uint, productID uint) error {
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	return s.repo.DecrementCartItemQuantity(cart.ID, productID)
}

func (s *cartServiceImpl) AddOneItem(userID uint, productID uint) error {
	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return err
	}

	item := domain.CartItem{
		ProductID: productID,
		Quantity:  1,
		Price:     product.Price,
	}

	return s.repo.AddOrUpdateCartItem(userID, item)
}

func (s *cartServiceImpl) Checkout(userID uint, req domain.CheckoutRequest) (*domain.Order, error) {
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(cart.CartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	order := domain.Order{
		UserID:        userID,
		AddressID:     req.ShippingAddressID,
		Status:        "pending",
		PaymentMethod: req.PaymentMethod,
		OrderItems:    make([]domain.OrderItem, len(cart.CartItems)),
	}

	for i, item := range cart.CartItems {
		order.OrderItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	order, err = s.orderService.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	err = s.repo.ClearCart(userID)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
