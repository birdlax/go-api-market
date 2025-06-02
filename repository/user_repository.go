package repository

import (
	"backend/domain"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
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

func (r *userRepositoryImpl) GetAll(page, limit int, sort, order string) ([]domain.User, int64, error) {
	var users []domain.User
	var totalItems int64

	offset := (page - 1) * limit

	// map sort field ป้องกัน SQL injection
	validSortFields := map[string]string{
		"id":        "id",
		"createdat": "created_at",
		"updatedat": "updated_at",
		"username":  "username",
		"email":     "email",
	}
	sortField := validSortFields[strings.ToLower(sort)]
	if sortField == "" {
		sortField = "created_at"
	}

	sortOrder := "ASC"
	if strings.ToLower(order) == "desc" {
		sortOrder = "DESC"
	}

	// นับทั้งหมดก่อน
	if err := r.db.Model(&domain.User{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// ดึงข้อมูลแบบมี offset / limit
	if err := r.db.
		Order(fmt.Sprintf("%s %s", sortField, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalItems, nil
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

func (r *userRepositoryImpl) SaveResetToken(userID uint, token string, expiresAt time.Time) error {
	prt := domain.ResetToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(&prt).Error
}

func (r *userRepositoryImpl) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepositoryImpl) FindUserIDByResetToken(token string) (uint, error) {
	var prt domain.ResetToken
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&prt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("invalid or expired token")
		}
		return 0, err
	}
	return prt.UserID, nil
}
func (r *userRepositoryImpl) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

// ลบ token หลังจากรีเซ็ตแล้ว
func (r *userRepositoryImpl) DeleteResetToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&domain.ResetToken{}).Error
}
