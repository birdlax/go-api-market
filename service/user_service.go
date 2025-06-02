package service

import (
	"backend/domain"
	"backend/utils"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	utils.Logger.Printf("🔍 [Service] Login attempt for email: %s", req.Email)

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		utils.Logger.Printf("❌ [Service] User not found: %s", req.Email)
		return nil, utils.NewAppError(404, "User not found")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.Logger.Printf("🔐 [Service] Invalid password for user: %s", req.Email)
		return nil, utils.NewAppError(401, "Invalid password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		utils.Logger.Printf("⚠️ [Service] Failed to generate token for user: %s, error: %v", req.Email, err)
		return nil, utils.NewAppError(500, "Failed to generate token")
	}

	utils.Logger.Printf("✅ [Service] User %s authenticated successfully", req.Email)

	user.Password = ""
	return &domain.LoginResponse{
		Token: token,
		User: &domain.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			Role:      user.Role,
		},
	}, nil
}

func (s *userService) GetByID(id uint) (*domain.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	var defaultAddr *domain.AddressResponse
	for _, a := range user.Addresses {
		if a.IsDefault {
			defaultAddr = &domain.AddressResponse{
				ID:           a.ID,
				FullName:     a.FullName,
				Phone:        a.Phone,
				AddressLine1: a.AddressLine1,
				AddressLine2: a.AddressLine2,
				City:         a.City,
				Province:     a.Province,
				ZipCode:      a.ZipCode,
				Country:      a.Country,
				IsDefault:    a.IsDefault,
			}
			break
		}
	}

	return &domain.UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Role:           user.Role,
		DefaultAddress: defaultAddr,
	}, nil
}

func (s *userService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *userService) GetAll(page, limit int, sort, order string) ([]domain.User, int64, error) {
	return s.repo.GetAll(page, limit, sort, order)
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

func (s *userService) FindByEmail(email string) (*domain.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) SendResetPasswordEmail(email string) error {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("email not found")
	}

	token := uuid.NewString()
	expiration := time.Now().Add(1 * time.Hour)

	if err := s.repo.SaveResetToken(user.ID, token, expiration); err != nil {
		return err
	}
	utils.Logger.Printf("Sending reset password email to: %s, token: %s", user.Email, token)
	// link := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token)
	return utils.SendResetPasswordEmail(user.Email, token)
}

func (s *userService) ResetPassword(token string, newPassword string) error {
	userID, err := s.repo.FindUserIDByResetToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(userID, hashedPassword); err != nil {
		return err
	}

	return s.repo.DeleteResetToken(token)
}
