package service

import (
	"backend/domain"
	"backend/utils"
	"errors"
	"fmt"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(email, password, role string, firstName *string, lastName *string) error {
	if email == "" || password == "" || role == "" {
		return errors.New("email, password, and role are required")
	}
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
		return fmt.Errorf("failed to get user: %w", err)
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}

	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}
