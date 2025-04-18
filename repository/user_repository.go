package repository

import (
	"backend/config"
	"backend/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(id uint) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	GetAll() ([]domain.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *domain.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	if err := config.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}
func (r *userRepository) Delete(id uint) error {
	if err := config.DB.Delete(&domain.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
func (r *userRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
