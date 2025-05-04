package handler

import (
	"backend/domain"
	"backend/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [CreateProduct] Start Create Product")
	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		utils.Logger.Printf("❌ [CreateProduct] Invalid request: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if product.CategoryID == 0 {
		utils.Logger.Printf("❌ [CreateProduct] Category ID is required")
		return c.Status(400).JSON(fiber.Map{"error": "Category ID is required"})
	}
	if err := h.service.CreateProduct(product); err != nil {
		utils.Logger.Printf("❌ [CreateProduct] Failed to Create Product: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [CreateProduct] Product created successfully")
	return c.JSON(fiber.Map{"message": "Product created successfully"})
}

func (h *ProductHandler) GetAllProduct(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [GetAllProduct] Start Get Product all")
	products, err := h.service.GetAllProduct()
	if err != nil {
		utils.Logger.Printf("❌ [GetAllProduct] Failed Get Product all: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [GetAllProduct] Get Product All successfully")
	return c.JSON(products)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [GetAllProducts] Start Get Product all")
	products, err := h.service.GetAllProducts()
	if err != nil {
		utils.Logger.Printf("❌ [GetAllProducts] Failed Get Product all: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [GetAllProduct] Get Product All successfully")
	return c.JSON(products)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	utils.Logger.Printf("🔄 [GetProductByID] Start Get product By ID: %s", id)
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Logger.Printf("[GetProductByID] Invalid ID format: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	product, err := h.service.GetProductByID(uint(parsedID))
	if err != nil {
		utils.Logger.Printf("❌ [GetProductByID] Get Product By ID not found: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [GetProductByID] Get Product By ID successfully")
	return c.JSON(product)
}

func (h *ProductHandler) GetproductByName(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [GetproductByName] Start Get Product By Name ")
	name := c.Params("name")
	product, err := h.service.GetProductByName(name)
	if err != nil {
		utils.Logger.Printf("❌ [GetproductByName] Get Product By Name not found: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [GetproductByName] Get Product By Name successfully")
	return c.JSON(product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	utils.Logger.Printf("🔄 [UpdateProduct] Start Update product: %s", id)
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Logger.Printf("❌ [UpdateProduct] Invalid ID format: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		utils.Logger.Printf("❌ [UpdateProduct] Invalid request: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	product.ID = uint(parsedID)

	if err := h.service.UpdateProduct(product); err != nil {
		utils.Logger.Printf("❌ [UpdateProduct] Failed to update product: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [UpdateProduct] Product updated successfully")
	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	utils.Logger.Printf("🔄 [Delete] Start Delete product: %s", id)
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Logger.Printf("❌ [Delete] Invalid ID format: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := h.service.Delete(uint(parsedID)); err != nil {
		utils.Logger.Printf("❌ [Delete] Failed to Delete product %d: %v", parsedID, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [Delete] Successfully Delete product ID: %d", parsedID)
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

// Category Handlers
func (h *ProductHandler) CreateCategory(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [CreateCategory] Start Create Category ")
	var category domain.Category
	if err := c.BodyParser(&category); err != nil {
		utils.Logger.Printf("❌ [CreateCategory] Invalid request: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.service.CreateCategory(category); err != nil {
		utils.Logger.Printf("❌ [Register] Failed Create Category %v: %v", category, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [CreateCategory] Create Category Successfully: %v", category)
	return c.JSON(fiber.Map{"message": "created Category successfully"})
}

func (h *ProductHandler) GetproductByCategory(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [GetproductByCategory] Start Get product By Category")
	category := c.Params("category")
	products, err := h.service.GetProductByCategory(category)
	if err != nil {
		utils.Logger.Printf("❌ [GetproductByCategory] Failed Get product By Category %v: %v", category, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [GetproductByCategory] Get product By Category Successfully: %v", category)
	return c.JSON(products)
}
