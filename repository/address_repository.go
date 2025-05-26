package repository

import (
	"backend/domain"
	"gorm.io/gorm"
)

type addressRepositoryImpl struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) domain.AddressRepository {
	return &addressRepositoryImpl{db: db}
}
func (r *addressRepositoryImpl) CreateAddress(address *domain.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepositoryImpl) GetAddressByID(id uint) (*domain.Address, error) {
	var address domain.Address
	if err := r.db.First(&address, id).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepositoryImpl) UpdateAddress(id uint, data domain.Address) error {
	return r.db.Model(&domain.Address{}).Where("id = ?", id).Updates(data).Error
}

func (r *addressRepositoryImpl) DeleteAddress(id uint) error {
	return r.db.Delete(&domain.Address{}, id).Error
}

func (r *addressRepositoryImpl) UnsetDefaultAddress(userID uint) error {
	return r.db.Model(&domain.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error
}

// คืนค่า Address ล่าสุดของผู้ใช้ (ใช้ ID ล่าสุด หรือ CreatedAt ล่าสุด)
func (r *addressRepositoryImpl) GetLatestAddressByUserID(userID uint) (*domain.Address, error) {
	var address domain.Address
	err := r.db.Where("user_id = ?", userID).Order("id DESC").First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}
func (r *addressRepositoryImpl) HasDefaultAddress(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Address{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *addressRepositoryImpl) GetAddressesByUserID(userID uint) ([]domain.Address, error) {
	var addresses []domain.Address
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, id DESC").Find(&addresses).Error
	return addresses, err
}
