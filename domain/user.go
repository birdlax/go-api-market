package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName *string `json:"first_name" gorm:"default:null"` // เปลี่ยนเป็น pointer
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
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type UpdateProfileRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
