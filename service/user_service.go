package service

import (
	"backend/domain"
	"backend/repository"
	"backend/utils"
	"errors"
)

type UserService interface {
	Register(email string, password string, role string, firstName, lastName *string) error
	Login(req domain.LoginRequest) (*domain.LoginResponse, error)

	GetByID(id uint) (*domain.UserResponse, error)
	Update(user *domain.User) error
	Delete(id uint) error
	GetAll() ([]domain.User, error)
	UpdatePassword(id uint, req domain.UpdatePasswordRequest) error
	UpdateProfile(id uint, req domain.UpdateProfileRequest) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// ปรับ service ให้รับ FirstName และ LastName ด้วย
func (s *userService) Register(email, password, role string, firstName *string, lastName *string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	existingUser, err := s.repo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return errors.New("user already exists")
	}

	user := &domain.User{
		Email:     email,
		Password:  hashedPassword,
		Role:      role,
		FirstName: firstName,
		LastName:  lastName,
	}

	return s.repo.Create(user)
}

func (s *userService) Login(req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("password is invalid")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("could not generate token")
	}
	user.Password = ""
	return &domain.LoginResponse{
		Token: token,
		User: &domain.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
		},
	}, nil
}

func (s *userService) GetByID(id uint) (*domain.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &domain.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}, nil
}

func (s *userService) Update(user *domain.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err := s.repo.Update(user); err != nil {
		return err
	}
	return nil
}

func (s *userService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *userService) GetAll() ([]domain.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) UpdatePassword(id uint, req domain.UpdatePasswordRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if err := s.repo.Update(user); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateProfile(id uint, req domain.UpdateProfileRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Role = req.Role

	if err := s.repo.Update(user); err != nil {
		return err
	}
	return nil
}
