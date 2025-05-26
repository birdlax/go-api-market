// models/address.go
package domain

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"not null"`
	Line1     string `json:"line1"`
	Line2     string `json:"line2"`
	City      string `json:"city"`
	Province  string `json:"province"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country" gorm:"default:'Thailand'"`
	IsDefault bool   `json:"is_default"`
}

type AddressRequest struct {
	Line1     string `json:"line1"`
	Line2     string `json:"line2"`
	City      string `json:"city"`
	Province  string `json:"province"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
	IsDefault bool   `json:"is_default"`
}

type AddressResponse struct {
	ID        uint   `json:"id"`
	Line1     string `json:"line1"`
	Line2     string `json:"line2"`
	City      string `json:"city"`
	Province  string `json:"province"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
	IsDefault bool   `json:"is_default"`
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
}
