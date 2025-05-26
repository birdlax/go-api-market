package repository

import (
	"backend/domain"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryImpl) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.
		Preload("Addresses", "is_default = ?", true).
		First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) Update(user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
func (r *userRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
func (r *userRepositoryImpl) GetAll() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepositoryImpl) GetDefaultAddress(userID uint) (*domain.Address, error) {
	var address domain.Address
	err := r.db.
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&address).Error

	if err != nil {
		return nil, err
	}
	return &address, nil
}
