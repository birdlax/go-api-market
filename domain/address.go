// models/address.go
package domain

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	UserID       uint   `json:"user_id" gorm:"not null"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	Province     string `json:"province"`
	ZipCode      string `json:"zip_code"`
	Country      string `json:"country" gorm:"default:'Thailand'"`
	IsDefault    bool   `json:"is_default"`
}

type AddressRequest struct {
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	Province     string `json:"province"`
	ZipCode      string `json:"zip_code"`
	Country      string `json:"country"`
	IsDefault    bool   `json:"is_default"`
}

type AddressResponse struct {
	ID           uint   `json:"id"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	Province     string `json:"province"`
	ZipCode      string `json:"zip_code"`
	Country      string `json:"country"`
	IsDefault    bool   `json:"is_default"`
}

type AddressRepository interface {
	CreateAddress(address *Address) error
	UpdateAddress(addressID uint, data Address) error
	DeleteAddress(addressID uint) error
	UnsetDefaultAddress(userID uint) error
	GetAddressByID(id uint) (*Address, error)
	GetLatestAddressByUserID(userID uint) (*Address, error)
	HasDefaultAddress(userID uint) (bool, error)
	GetAddressesByUserID(userID uint) ([]Address, error)
}

type AddressService interface {
	CreateAddress(userID uint, address *Address) error
	UpdateAddress(addressID uint, req AddressRequest) error
	DeleteAddress(addressID uint) error
	GetAddressesByUserID(userID uint) ([]AddressResponse, error)
	GetAddressByID(id uint) (*Address, error)
	SwitchDefaultAddress(userID uint, addressID uint) error
}
