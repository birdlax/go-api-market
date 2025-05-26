package service

import (
	"backend/domain"
	"fmt"
	"time"
)

type orderServiceImpl struct {
	repo     domain.OrderRepository
	cartRepo domain.CartRepository
}

func NewOrderService(repo domain.OrderRepository, cartRepo domain.CartRepository) domain.OrderService {
	return &orderServiceImpl{repo: repo, cartRepo: cartRepo}
}

func (s *orderServiceImpl) CreateOrder(order domain.Order) (domain.Order, error) {
	existingOrder, err := s.repo.GetPendingOrderByUserID(order.UserID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to check existing order: %v", err)
	}

	tx := s.repo.BeginTx()

	// ถ้ามีคำสั่งซื้อที่ค้างอยู่ ให้ลบสินค้าเก่าแล้วเพิ่มสินค้าใหม่
	if existingOrder.ID != 0 {
		if err := s.repo.DeleteOrderItemsByOrderID(tx, existingOrder.ID); err != nil {
			tx.Rollback()
			return domain.Order{}, fmt.Errorf("failed to clear old order items: %v", err)
		}

		var totalPrice float64
		var newOrderItems []domain.OrderItem

		for _, item := range order.OrderItems {
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
			item.OrderID = existingOrder.ID
			totalPrice += item.Price
			newOrderItems = append(newOrderItems, item)
		}

		if err := s.repo.CreateOrderItems(tx, newOrderItems); err != nil {
			tx.Rollback()
			return domain.Order{}, fmt.Errorf("failed to create new order items: %v", err)
		}

		existingOrder.TotalPrice = totalPrice
		existingOrder.Status = "pending"

		if err := s.repo.UpdateOrderWithTx(tx, &existingOrder); err != nil {
			tx.Rollback()
			return domain.Order{}, fmt.Errorf("failed to update order: %v", err)
		}

		if err := tx.Commit().Error; err != nil {
			return domain.Order{}, err
		}

		return existingOrder, nil
	}

	// ถ้าไม่มีคำสั่งซื้อที่ยังค้างอยู่ ให้สร้างใหม่
	var totalPrice float64
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
		order.OrderItems[i].Price = item.Price
		totalPrice += item.Price
	}

	order.TotalPrice = totalPrice
	order.Status = "pending"

	if err := s.repo.CreateWithTx(tx, &order); err != nil {
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

func (s *orderServiceImpl) GetUnpaidOrdersByUserID(userID uint) ([]domain.Order, error) {
	return s.repo.GetOrdersByUserIDAndStatus(userID, "pending")
}
func (s *orderServiceImpl) GetOrdersByUserIDAndStatus(userID uint, status string) ([]domain.Order, error) {
	return s.repo.GetOrdersByUserIDAndStatus(userID, status)
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

	if len(updated.OrderItems) == 0 {
		return domain.Order{}, fmt.Errorf("order must contain at least one item")
	}

	var totalPrice float64
	for _, item := range updated.OrderItems {
		product, err := s.repo.GetProductByID(item.ProductID)
		if err != nil {
			return domain.Order{}, fmt.Errorf("product %d not found", item.ProductID)
		}

		// คำนวณราคาและตรวจสอบสต๊อกสินค้า
		if product.Quantity < item.Quantity {
			return domain.Order{}, fmt.Errorf("not enough stock for product %s", product.Name)
		}

		item.Price = product.Price * float64(item.Quantity)
		totalPrice += item.Price
	}

	order.TotalPrice = totalPrice
	order.OrderItems = updated.OrderItems

	if updated.Status != "" {
		order.Status = updated.Status
	}

	return s.repo.UpdateOrder(*order)
}

func (s *orderServiceImpl) DeleteOrder(id uint) error {
	return s.repo.DeleteOrder(id)
}

func (s *orderServiceImpl) MarkOrderAsPaidByUserID(userID uint) error {
	// ดึงคำสั่งซื้อที่ยังไม่ได้ชำระของผู้ใช้
	order, err := s.repo.GetPendingOrderByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get pending order: %w", err)
	}

	// ตรวจสอบว่าพบคำสั่งซื้อหรือไม่
	if order.ID == 0 {
		return fmt.Errorf("no pending order found for user")
	}

	// ตรวจสอบสถานะ
	if order.Status == "paid" {
		return fmt.Errorf("order already paid")
	}

	now := time.Now()
	order.Status = "paid"
	order.PaidAt = &now

	// อัพเดตสถานะคำสั่งซื้อ
	if _, err := s.repo.UpdateOrder(order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// ล้างตะกร้าสินค้าของผู้ใช้
	if err := s.cartRepo.ClearCart(order.UserID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

func (s *orderServiceImpl) CancelOrderByUserID(userID uint) error {
	// ดึงคำสั่งซื้อที่ยังไม่ชำระของผู้ใช้
	order, err := s.repo.GetPendingOrderByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get pending order: %w", err)
	}

	// ตรวจสอบว่าพบคำสั่งซื้อหรือไม่
	if order.ID == 0 {
		return fmt.Errorf("no pending order found for user")
	}

	// ตรวจสอบสถานะ
	if order.Status == "paid" {
		return fmt.Errorf("cannot cancel a paid order")
	}
	if order.Status == "canceled" {
		return fmt.Errorf("order already canceled")
	}

	order.Status = "canceled"

	// อัพเดตสถานะคำสั่งซื้อ
	if _, err := s.repo.UpdateOrder(order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// ล้างตะกร้าสินค้าหากต้องการ
	if err := s.cartRepo.ClearCart(order.UserID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

// Service
func (s *orderServiceImpl) GetRevenueByCategory(status string) ([]domain.RevenueResult, error) {
	results, err := s.repo.GetRevenueByCategory(status)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue by category: %w", err)
	}
	return results, nil
}
