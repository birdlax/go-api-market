package handler

import (
	"backend/domain"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	service domain.CartService
}

func NewCartHandler(service domain.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) AddItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var input domain.CartItemInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	log.Println("AddItem input:", input)
	product, err := h.service.GetProductByID(input.ProductID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	item := domain.CartItem{
		CartID:    userID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
		Price:     product.Price, // ดึงราคาจากฐานข้อมูล
	}

	if err := h.service.AddItem(userID, item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Item added to cart"})
}

// RemoveItem สำหรับลบสินค้าจากตะกร้า
func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	productID := c.Params("product_id") // parse uint if needed

	// convert productID to uint
	id, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	if err := h.service.RemoveItem(userID, uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Item removed from cart"})
}

// GetCart สำหรับดึงข้อมูลตะกร้าสินค้าทั้งหมดของ user
func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	cart, err := h.service.GetCart(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cart)
}

func (h *CartHandler) Checkout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if err := h.service.Checkout(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Checkout successful"})
}
