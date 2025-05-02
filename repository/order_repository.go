package repository

import (
	"backend/domain"
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

func (r *orderRepositoryImpl) CreateWithTx(tx *gorm.DB, order domain.Order) error {
	return tx.Create(&order).Error
}

func (r *orderRepositoryImpl) GetProductByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *orderRepositoryImpl) GetAllOrders() ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Preload("OrderItems").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepositoryImpl) UpdateProductStock(tx *gorm.DB, product *domain.Product) error {
	return tx.Save(product).Error
}

func (r *orderRepositoryImpl) GetOrderByID(id uint) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.Preload("OrderItems").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepositoryImpl) UpdateOrder(order domain.Order) (domain.Order, error) {
	// 🔄 ลบ OrderItems เดิมก่อน
	if err := r.db.Where("order_id = ?", order.ID).Delete(&domain.OrderItem{}).Error; err != nil {
		return order, err
	}

	// 💾 Save พร้อม OrderItems ใหม่
	err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&order).Error
	return order, err
}

func (r *orderRepositoryImpl) DeleteOrder(id uint) error {
	return r.db.Delete(&domain.Order{}, id).Error
}
