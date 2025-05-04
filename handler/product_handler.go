package handler

import (
	"backend/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if product.CategoryID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Category ID is required"})
	}
	if err := h.service.CreateProduct(product); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Product created successfully"})
}

func (h *ProductHandler) GetAllProduct(c *fiber.Ctx) error {
	products, err := h.service.GetAllProduct()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	product, err := h.service.GetProductByID(uint(parsedID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(product)
}

func (h *ProductHandler) GetproductByName(c *fiber.Ctx) error {
	name := c.Params("name")
	product, err := h.service.GetProductByName(name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	product.ID = uint(parsedID)

	if err := h.service.UpdateProduct(product); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := h.service.Delete(uint(parsedID)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

// Category Handlers
func (h *ProductHandler) CreateCategory(c *fiber.Ctx) error {
	var category domain.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.service.CreateCategory(category); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Category created successfully"})
}

func (h *ProductHandler) GetproductByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	products, err := h.service.GetProductByCategory(category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}
