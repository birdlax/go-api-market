package handler

import (
	"backend/domain"
	"backend/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
		utils.Logger.Printf("AddItem: invalid input - error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	utils.Logger.Printf("AddItem input - userID: %d, productID: %d, quantity: %d", userID, input.ProductID, input.Quantity)

	product, err := h.service.GetProductByID(input.ProductID)
	if err != nil {
		utils.Logger.Printf("AddItem: product not found - productID: %d, error: %v", input.ProductID, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	item := domain.CartItem{
		CartID:    userID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
		Price:     product.Price,
	}

	if err := h.service.AddItem(userID, item); err != nil {
		utils.Logger.Printf("AddItem: failed to add item - userID: %d, error: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Printf("AddItem: success - userID: %d, productID: %d", userID, input.ProductID)
	return c.JSON(fiber.Map{"message": "Item added to cart"})
}

func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	productID := c.Params("product_id")

	id, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		utils.Logger.Printf("RemoveItem: invalid product ID - productID: %s, error: %v", productID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	if err := h.service.RemoveItem(userID, uint(id)); err != nil {
		utils.Logger.Printf("RemoveItem: failed to remove item - userID: %d, productID: %d, error: %v", userID, id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Printf("RemoveItem: success - userID: %d, productID: %d", userID, id)
	return c.JSON(fiber.Map{"message": "Item removed from cart"})
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	cart, err := h.service.GetCart(userID)
	if err != nil {
		utils.Logger.Printf("GetCart: failed to get cart - userID: %d, error: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Printf("GetCart: success - userID: %d", userID)
	return c.JSON(cart)
}

func (h *CartHandler) Checkout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	if err := h.service.Checkout(userID); err != nil {
		utils.Logger.Printf("Checkout: failed - userID: %d, error: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Printf("Checkout: success - userID: %d", userID)
	return c.JSON(fiber.Map{"message": "Checkout successful"})
}
