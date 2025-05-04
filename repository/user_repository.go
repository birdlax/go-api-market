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
	if err := r.db.First(&user, id).Error; err != nil {
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
