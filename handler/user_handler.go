package handler

import (
	"backend/config"
	"backend/domain"
	"backend/utils"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [Register] Start user registration")

	var req domain.User
	if err := c.BodyParser(&req); err != nil {
		utils.Logger.Printf("❌ [Register] Invalid request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.FirstName == nil || req.LastName == nil || *req.FirstName == "" || *req.LastName == "" {
		utils.Logger.Printf("❌ [Register] First name and last name are required")
		return c.Status(400).JSON(fiber.Map{"error": "First name and last name are required"})
	}

	if req.Role == "" {
		req.Role = "user"
		utils.Logger.Println("ℹ️ [Register] Role not provided, defaulting to 'user'")
	}

	utils.Logger.Printf("📥 [Register] Registering user: %s", req.Email)

	err := h.service.Register(req.Email, req.Password, req.Role, req.FirstName, req.LastName)
	if err != nil {
		utils.Logger.Printf("❌ [Register] Failed to register user %s: %v", req.Email, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Printf("✅ [Register] User %s registered successfully", req.Email)
	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [Handler] Start login process")

	var req domain.User
	if err := c.BodyParser(&req); err != nil {
		utils.Logger.Printf("❌ [Handler] Invalid request format: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	utils.Logger.Printf("📥 [Handler] Login request received for email: %s", req.Email)

	user, err := h.service.Login(domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if appErr, ok := utils.AsAppError(err); ok {
			utils.Logger.Printf("❗ [Handler] Login error: %s - %s", req.Email, appErr.Message)
			return c.Status(appErr.Code).JSON(fiber.Map{"error": appErr.Message})
		}
		utils.Logger.Printf("🔥 [Handler] Unexpected error for email %s: %v", req.Email, err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	isSecureConnection := c.Protocol() == "https"
	c.Cookie(&fiber.Cookie{
		Name:     config.JwtCookieName,
		Value:    user.Token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
		Secure:   isSecureConnection,
		SameSite: "Lax",
		Path:     "/",
	})

	utils.Logger.Printf("🎉 [Handler] User %s logged in successfully", req.Email)
	return c.JSON(user)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	utils.Logger.Printf("🔄 [GetByID] Start get user by ID: %s", id)

	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Logger.Printf("[GetByID] Invalid ID format: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	utils.Logger.Printf("📥 [GetByID] Getting user with ID: %d", parsedID)

	user, err := h.service.GetByID(uint(parsedID))
	if err != nil {
		utils.Logger.Printf("❌ [GetByID] User not found: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	utils.Logger.Printf("✅ [GetByID] User ID %d retrieved successfully", parsedID)
	return c.JSON(user)
}

func (h *UserHandler) UpdateProfilebyId(c *fiber.Ctx) error {
	id := c.Params("id")
	utils.Logger.Printf("🔄 [UpdateProfilebyId] Updating profile for user ID %s", id)
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Logger.Printf("❌ [UpdateProfilebyId] Invalid ID format: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var req domain.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Logger.Printf("❌ [UpdateProfilebyId] Invalid request: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err = h.service.UpdateProfile(uint(parsedID), req)
	if err != nil {
		utils.Logger.Printf("❌ [UpdateProfilebyId] Failed to update user %d: %v", parsedID, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [UpdateProfilebyId] Successfully updated user ID: %d", parsedID)
	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	utils.Logger.Printf("🔄 [Delete] Start to delete user ID: %s", id)

	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Logger.Printf("❌ [Delete] Invalid ID format: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	err = h.service.Delete(uint(parsedID))
	if err != nil {
		utils.Logger.Printf("❌ [Delete] Failed to Delete user %d: %v", parsedID, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [Delete] Successfully Delete user ID: %d", parsedID)
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [GetAll] Start retrieving all users")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	users, totalItems, err := h.service.GetAll(page, limit, sort, order)
	if err != nil {
		utils.Logger.Printf("❌ [GetAll] Failed to get users: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return c.JSON(fiber.Map{
		"current_page": page,
		"items":        users,
		"per_page":     limit,
		"total_items":  totalItems,
		"total_pages":  totalPages,
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [Logout] Start logging out user")

	c.Cookie(&fiber.Cookie{
		Name:     "JWT",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	utils.Logger.Println("✅ [Logout] User logged out successfully")
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	utils.Logger.Printf("🔄 [GetCurrentUser] Start to GetCurrentUser user ID: %s", userID)
	if userID == nil {
		utils.Logger.Printf("❌ [GetCurrentUser] User ID not found in context")
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	user, err := h.service.GetByID(userID.(uint))
	if err != nil {
		utils.Logger.Printf("❌ [GetCurrentUser] Failed to get user with ID %v: %v", userID, err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	utils.Logger.Printf("✅ [GetCurrentUser] Successfully retrieved user ID: %v", userID)
	return c.JSON(user)
}

func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	utils.Logger.Printf("🔄 [UpdatePassword] Start to UpdatePassword user ID: %s", userID)
	if userID == nil {
		utils.Logger.Printf("❌ [UpdatePassword] User ID not found in context")
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req domain.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Logger.Printf("❌ [UpdatePassword] Invalid request %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.service.UpdatePassword(userID.(uint), req)
	if err != nil {
		utils.Logger.Printf("❌ [UpdatePassword] Failed to update password for user ID %v: %v", userID, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Printf("✅ [UpdatePassword] Successfully updated password for user ID %v", userID)
	return c.JSON(fiber.Map{"message": "Password updated successfully"})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		utils.Logger.Printf("❌ [UpdateProfile] User ID not found in context")
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	utils.Logger.Printf("🔄 [UpdateProfile] Start updating profile for user ID: %v", userID)

	var req domain.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Logger.Printf("❌ [UpdateProfile] Invalid request : %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if *req.FirstName == "" || *req.LastName == "" {
		utils.Logger.Printf("❌ [UpdateProfile] Missing required fields for user ID %v", userID)
		return c.Status(400).JSON(fiber.Map{"error": "First name and Last name are required"})
	}

	err := h.service.UpdateProfile(userID.(uint), req)
	if err != nil {
		utils.Logger.Printf("❌ [UpdateProfile] Failed to update profile for user ID %v: %v", userID, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [UpdateProfile] Successfully updated profile for user ID %v", userID)
	return c.JSON(fiber.Map{"message": "Profile updated successfully"})
}

func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.service.SendResetPasswordEmail(req.Email); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Reset link sent"})
}

func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Token == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token and new password are required",
		})
	}

	err := h.service.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}

func (h *UserHandler) GetHello(c *fiber.Ctx) error {
	return c.SendString("Hello, Backend Golang v1!")
}
