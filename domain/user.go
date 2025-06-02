package domain

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	FirstName *string   `json:"first_name" gorm:"default:null"`
	LastName  *string   `json:"last_name" gorm:"default:null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Role      string    `json:"role"`
	Addresses []Address `gorm:"foreignKey:UserID"`
}

type UserResponse struct {
	ID             uint             `json:"id"`
	Email          string           `json:"email"`
	FirstName      *string          `json:"first_name"`
	LastName       *string          `json:"last_name"`
	Role           string           `json:"role"`
	DefaultAddress *AddressResponse `json:"default_address,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error

	SaveResetToken(userID uint, token string, expiresAt time.Time) error
	FindByEmail(email string) (*User, error)
	FindUserIDByResetToken(token string) (uint, error)
	UpdatePassword(userID uint, hashedPassword string) error
	DeleteResetToken(token string) error
	GetAll(page, limit int, sort, order string) ([]User, int64, error)
}

type UserService interface {
	Register(email string, password string, role string, firstName, lastName *string) error
	Login(req LoginRequest) (*LoginResponse, error)
	GetByID(id uint) (*UserResponse, error)
	Delete(id uint) error
	UpdatePassword(id uint, req UpdatePasswordRequest) error
	UpdateProfile(id uint, req UpdateProfileRequest) error

	FindByEmail(email string) (*User, error)
	SendResetPasswordEmail(email string) error
	ResetPassword(token string, newPassword string) error
	GetAll(page, limit int, sort, order string) ([]User, int64, error)
}
