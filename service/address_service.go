package service

import (
	"backend/domain"
)

type addressServiceImpl struct {
	repo domain.AddressRepository
}

func NewAddressService(repo domain.AddressRepository) domain.AddressService {
	return &addressServiceImpl{repo: repo}
}

func (s *addressServiceImpl) CreateAddress(userID uint, address *domain.Address) error {
	address.UserID = userID

	if address.IsDefault {
		// ถ้ามีการกำหนดเป็น default -> unset default ตัวอื่นก่อน
		if err := s.repo.UnsetDefaultAddress(userID); err != nil {
			return err
		}
	} else {
		// ถ้าไม่ได้ตั้งให้เป็น default ต้องเช็คว่ามี default อยู่รึยัง
		hasDefault, err := s.repo.HasDefaultAddress(userID)
		if err != nil {
			return err
		}
		if !hasDefault {
			// ถ้ายังไม่มี default ตัวไหนเลย -> ให้ตัวนี้เป็น default
			address.IsDefault = true
		}
	}

	return s.repo.CreateAddress(address)
}

func (s *addressServiceImpl) UpdateAddress(addressID uint, req domain.AddressRequest) error {
	address, err := s.repo.GetAddressByID(addressID)
	if err != nil {
		return err
	}

	if req.IsDefault {
		if err := s.repo.UnsetDefaultAddress(address.UserID); err != nil {
			return err
		}
	}

	updated := domain.Address{
		Line1:     req.Line1,
		Line2:     req.Line2,
		City:      req.City,
		Province:  req.Province,
		ZipCode:   req.ZipCode,
		Country:   req.Country,
		IsDefault: req.IsDefault,
	}

	return s.repo.UpdateAddress(addressID, updated)
}

func (s *addressServiceImpl) DeleteAddress(addressID uint) error {
	address, err := s.repo.GetAddressByID(addressID)
	if err != nil {
		return err
	}

	// ลบ address
	if err := s.repo.DeleteAddress(addressID); err != nil {
		return err
	}

	// ถ้า address ที่ลบเป็น default ให้ตั้งตัวอื่นเป็น default แทน (ถ้ามี)
	if address.IsDefault {
		// หา address ใหม่ของ user
		newDefault, err := s.repo.GetLatestAddressByUserID(address.UserID)
		if err == nil && newDefault != nil {
			_ = s.repo.UpdateAddress(newDefault.ID, domain.Address{IsDefault: true})
		}
	}

	return nil
}

func (s *addressServiceImpl) GetAddressesByUserID(userID uint) ([]domain.AddressResponse, error) {
	addresses, err := s.repo.GetAddressesByUserID(userID)
	if err != nil {
		return nil, err
	}

	var res []domain.AddressResponse
	for _, addr := range addresses {
		res = append(res, domain.AddressResponse{
			ID:        addr.ID,
			Line1:     addr.Line1,
			Line2:     addr.Line2,
			City:      addr.City,
			Province:  addr.Province,
			ZipCode:   addr.ZipCode,
			Country:   addr.Country,
			IsDefault: addr.IsDefault,
		})
	}
	return res, nil
}
