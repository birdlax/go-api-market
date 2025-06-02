package repository

import (
	"backend/domain"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepositoryImpl{db: db}
}

func (r *orderRepositoryImpl) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *orderRepositoryImpl) CreateWithTx(tx *gorm.DB, order *domain.Order) error {
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback() // ยกเลิก transaction ถ้าเกิด error
		return fmt.Errorf("failed to update order: %v", err)
	}

	// ไม่มี error แสดงว่าอัพเดทคำสั่งซื้อสำเร็จ
	return nil
}

func (r *orderRepositoryImpl) GetOrdersByUserIDAndStatus(userID uint, status string) ([]domain.Order, error) {
	var orders []domain.Order

	query := r.db.Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if status != "" {
		query = query.Where("status = ?", status)

	}

	err := query.Find(&orders).Error
	return orders, err
}

func (r *orderRepositoryImpl) GetProductByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *orderRepositoryImpl) GetAllOrders(page, limit int, sort, order string) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var totalItems int64

	offset := (page - 1) * limit

	validSortFields := map[string]string{
		"createdat": "created_at",
		"updatedat": "updated_at",
		"id":        "id",
	}

	sortField := validSortFields[strings.ToLower(sort)]
	if sortField == "" {
		sortField = "created_at"
	}

	sortOrder := "ASC"
	if strings.ToLower(order) == "desc" {
		sortOrder = "DESC"
	}

	// นับจำนวนรวม
	if err := r.db.Model(&domain.Order{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.
		Preload("OrderItems").
		Order(fmt.Sprintf("%s %s", sortField, sortOrder)).
		Limit(limit).
		Offset(offset)

	if err := query.Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, totalItems, nil
}

func (r *orderRepositoryImpl) UpdateProductStock(tx *gorm.DB, product *domain.Product) error {
	return tx.Save(product).Error
}

func (r *orderRepositoryImpl) GetOrderByID(id uint) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("OrderItems.Product").
		Preload("Address").
		First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepositoryImpl) UpdateOrder(order domain.Order) (domain.Order, error) {
	if err := r.db.Where("order_id = ?", order.ID).Delete(&domain.OrderItem{}).Error; err != nil {
		return order, err
	}

	err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&order).Error
	return order, err
}

func (r *orderRepositoryImpl) DeleteOrder(id uint) error {
	tx := r.db.Begin()
	if err := tx.Where("order_id = ?", id).Delete(&domain.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&domain.Order{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *orderRepositoryImpl) GetPendingOrderByUserID(userID uint) (domain.Order, error) {
	var order domain.Order
	err := r.db.Where("user_id = ? AND status = ?", userID, "pending").Order("created_at DESC").First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Order{}, nil // ไม่พบ order ที่ค้างอยู่ = OK
		}
		return domain.Order{}, err // error อื่น ๆ
	}
	return order, nil // เจอ order ที่ค้างอยู่
}

func (r *orderRepositoryImpl) UpdateOrderWithTx(tx *gorm.DB, order *domain.Order) error {
	// ใช้ tx.Save เพื่ออัพเดตคำสั่งซื้อที่มีอยู่
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback() // หากเกิดข้อผิดพลาด ให้ทำการ Rollback
		return err
	}

	return nil
}
func (r *orderRepositoryImpl) DeleteOrderItemsByOrderID(tx *gorm.DB, orderID uint) error {
	return tx.Where("order_id = ?", orderID).Delete(&domain.OrderItem{}).Error
}

func (r *orderRepositoryImpl) CreateOrderItems(tx *gorm.DB, items []domain.OrderItem) error {
	return tx.Create(&items).Error
}
