package handler

import (
	"backend/domain"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type AddressHandler struct {
	service domain.AddressService
}

func NewAddressHandler(service domain.AddressService) *AddressHandler {
	return &AddressHandler{service: service}
}
func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var address domain.Address
	if err := c.BodyParser(&address); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.CreateAddress(userID, &address); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(address)
}

// func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
// 	id, _ := strconv.Atoi(c.Params("id"))
// 	var req domain.AddressRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
// 	}

// 	if err := h.service.UpdateAddress(uint(id), req); err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{"message": "Address updated"})
// }

// func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
// 	id, _ := strconv.Atoi(c.Params("id"))

// 	if err := h.service.DeleteAddress(uint(id)); err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{"message": "Address deleted"})
// }

func (h *AddressHandler) GetAddresses(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	addresses, err := h.service.GetAddressesByUserID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(addresses)
}

func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))

	// ตรวจสอบว่า address เป็นของ user หรือไม่
	address, err := h.service.GetAddressByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Address not found"})
	}
	if address.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req domain.AddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.service.UpdateAddress(uint(id), req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Address updated"})
}

func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))

	// ตรวจสอบว่า address เป็นของ user หรือไม่
	address, err := h.service.GetAddressByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Address not found"})
	}
	if address.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.service.DeleteAddress(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Address deleted"})
}
func (h *AddressHandler) GetAddressByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	address, err := h.service.GetAddressByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Address not found"})
	}

	return c.JSON(address)
}

func (h *AddressHandler) SwitchDefault(c *fiber.Ctx) error {
	addressID, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(uint)

	err := h.service.SwitchDefaultAddress(userID, uint(addressID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Default address updated"})
}
