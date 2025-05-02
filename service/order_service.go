package service

import (
	"backend/domain"
	"fmt"
)

type OrderService interface {
	CreateOrder(order domain.Order) (domain.Order, error)
	GetAllOrders() ([]domain.Order, error)
	GetOrderByID(id uint) (domain.Order, error)
	UpdateOrder(id uint, updated domain.Order) (domain.Order, error)
	DeleteOrder(id uint) error
}

type orderServiceImpl struct {
	repo domain.OrderRepository
}

func NewOrderService(repo domain.OrderRepository) OrderService {
	return &orderServiceImpl{repo: repo}
}

func (s *orderServiceImpl) CreateOrder(order domain.Order) (domain.Order, error) {
	var totalPrice float64

	tx := s.repo.BeginTx()

	for i, item := range order.OrderItems {
		if item.Quantity <= 0 {
			tx.Rollback()
			return domain.Order{}, fmt.Errorf("quantity must be greater than 0 for product %d", item.ProductID)
		}

		product, err := s.repo.GetProductByID(item.ProductID)
		if err != nil {
			tx.Rollback()
			return domain.Order{}, fmt.Errorf("product %d not found", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			tx.Rollback()
			return domain.Order{}, fmt.Errorf("not enough stock for product %s", product.Name)
		}

		product.Quantity -= item.Quantity
		if err := s.repo.UpdateProductStock(tx, product); err != nil {
			tx.Rollback()
			return domain.Order{}, fmt.Errorf("failed to update stock for product %s", product.Name)
		}

		item.Price = product.Price * float64(item.Quantity)
		totalPrice += item.Price

		order.OrderItems[i].Price = item.Price
	}

	order.TotalPrice = totalPrice
	order.Status = "pending"

	if err := s.repo.CreateWithTx(tx, order); err != nil {
		tx.Rollback()
		return domain.Order{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (s *orderServiceImpl) GetAllOrders() ([]domain.Order, error) {
	return s.repo.GetAllOrders()

}

func (s *orderServiceImpl) GetOrderByID(id uint) (domain.Order, error) {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return domain.Order{}, err
	}
	return *order, nil
}

func (s *orderServiceImpl) UpdateOrder(id uint, updated domain.Order) (domain.Order, error) {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return domain.Order{}, err
	}
	if order.Status == "paid" || order.Status == "shipped" {
		return domain.Order{}, fmt.Errorf("cannot update order after payment or shipment")
	}
	if updated.TotalPrice <= 0 {
		return domain.Order{}, fmt.Errorf("total price must be greater than 0")
	}
	if len(updated.OrderItems) == 0 {
		return domain.Order{}, fmt.Errorf("order must contain at least one item")
	}
	if updated.Status != "" {
		order.Status = updated.Status
	}
	order.TotalPrice = updated.TotalPrice
	order.OrderItems = updated.OrderItems

	return s.repo.UpdateOrder(*order)
}

func (s *orderServiceImpl) DeleteOrder(id uint) error {
	return s.repo.DeleteOrder(id)
}
