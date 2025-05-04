package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName *string `json:"first_name" gorm:"default:null"`
	LastName  *string `json:"last_name" gorm:"default:null"`
	Email     string  `gorm:"unique;not null"`
	Password  string  `gorm:"not null"`
	Role      string  `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}
type UserResponse struct {
	ID        uint    `json:"id"`
	Email     string  `json:"email"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Role      string  `json:"role"`
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

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	GetAll() ([]User, error)
}

type UserService interface {
	Register(email string, password string, role string, firstName, lastName *string) error
	Login(req LoginRequest) (*LoginResponse, error)

	GetByID(id uint) (*UserResponse, error)
	Delete(id uint) error
	GetAll() ([]User, error)
	UpdatePassword(id uint, req UpdatePasswordRequest) error
	UpdateProfile(id uint, req UpdateProfileRequest) error
}
