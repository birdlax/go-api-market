package handler

import (
	"backend/domain"
	"backend/utils"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateMultipleProducts(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid multipart form"})
	}

	products := []*domain.Product{}

	names := form.Value["name"]
	descriptions := form.Value["description"]
	prices := form.Value["price"]
	quantities := form.Value["quantity"]
	categoryIDs := form.Value["category_id"]
	files := form.File["images"] // ใช้ชื่อฟิลด์ "images" และอัปโหลดหลายไฟล์

	for i := 0; i < len(names); i++ {
		categoryID, _ := strconv.Atoi(categoryIDs[i])
		price, _ := strconv.ParseFloat(prices[i], 64)
		quantity, _ := strconv.Atoi(quantities[i])

		product := &domain.Product{
			Name:        names[i],
			Description: descriptions[i],
			CategoryID:  uint(categoryID),
			Price:       price,
			Quantity:    quantity,
			Images:      []domain.ProductImage{},
		}

		uploadDir := fmt.Sprintf("./uploads/%d", categoryID)
		os.MkdirAll(uploadDir, os.ModePerm)

		// กรองเฉพาะไฟล์ภาพชุดของ product นี้ (เช่น ตาม index หรือชื่อฟิลด์)
		for _, file := range files {
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			filePath := fmt.Sprintf("%s/%s", uploadDir, filename)

			if err := c.SaveFile(file, filePath); err != nil {
				log.Println("save failed:", err)
				continue
			}

			// เพิ่มเข้า images
			product.Images = append(product.Images, domain.ProductImage{
				Path: filePath,
			})
		}

		products = append(products, product)
	}

	created, skipped, err := h.service.CreateProducts(products)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":          "Products created",
		"created_count":    len(created),
		"skipped_count":    len(skipped),
		"skipped_products": skipped,
	})
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
	searchQuery := c.Query("q", "")
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}

	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	minPriceStr := c.Query("min_price", "")
	maxPriceStr := c.Query("max_price", "")

	var minPrice, maxPrice float64
	if minPriceStr != "" {
		minPrice, _ = strconv.ParseFloat(minPriceStr, 64)
	}
	if maxPriceStr != "" {
		maxPrice, _ = strconv.ParseFloat(maxPriceStr, 64)
	}

	products, totalItems, err := h.service.GetAllProducts(page, limit, sort, order, minPrice, maxPrice, searchQuery)
	if err != nil {
		utils.Logger.Printf("❌ [GetAllProducts] Failed Get Product all: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return c.JSON(fiber.Map{
		"items":        products,
		"total_items":  totalItems,
		"total_pages":  totalPages,
		"current_page": page,
		"per_page":     limit,
	})
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
		utils.Logger.Printf("❌ [GetProductByID] Failed to get product with ID %d: %v", parsedID, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [GetProductByID] Get Product By ID successfully")
	return c.JSON(product)
}

func (h *ProductHandler) GetProductByName(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [GetproductByName] Start Get Product By Name ")
	name := c.Params("name")
	product, err := h.service.GetProductByName(name)
	if err != nil {
		utils.Logger.Printf("❌ [GetProductByName] Failed to get product with name '%s': %v", name, err)
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
		utils.Logger.Printf("❌ [UpdateProduct] Invalid request body:: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	product.ID = uint(parsedID)

	if err := h.service.UpdateProduct(product); err != nil {
		utils.Logger.Printf("❌ [UpdateProduct] Failed to update product ID %d: %v", parsedID, err)
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
		utils.Logger.Printf("❌ [CreateCategory] Invalid request body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.service.CreateCategory(category); err != nil {
		utils.Logger.Printf("❌ [CreateCategory] Failed to create category: %+v, error: %v", category, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [CreateCategory] Create Category Successfully: %v", category)
	return c.JSON(fiber.Map{"message": "created Category successfully"})
}

func (h *ProductHandler) GetProductByCategory(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [GetProductByCategory] Start Get product By Category")

	category := c.Params("category")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}

	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	minPriceStr := c.Query("min_price", "")
	maxPriceStr := c.Query("max_price", "")

	var minPrice, maxPrice float64
	if minPriceStr != "" {
		minPrice, _ = strconv.ParseFloat(minPriceStr, 64)
	}
	if maxPriceStr != "" {
		maxPrice, _ = strconv.ParseFloat(maxPriceStr, 64)
	}

	products, totalItems, err := h.service.GetProductByCategory(category, page, limit, sort, order, minPrice, maxPrice)
	if err != nil {
		utils.Logger.Printf("❌ [GetProductByCategory] Failed: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return c.JSON(fiber.Map{
		"items":        products,
		"total_items":  totalItems,
		"total_pages":  totalPages,
		"current_page": page,
		"per_page":     limit,
	})
}

func (h *ProductHandler) CreateMultipleCategories(c *fiber.Ctx) error {
	var categories []domain.Category // สมมุติใช้ GORM model

	if err := c.BodyParser(&categories); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	for _, cat := range categories {
		if err := h.service.CreateCategory(cat); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to create category: %s", cat.Name),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Categories created successfully",
	})
}

func (h *ProductHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	return c.Status(fiber.StatusOK).JSON(categories)
}

func (h *ProductHandler) GetNewArrivals(c *fiber.Ctx) error {
	utils.Logger.Println("🆕 [GetNewArrivals] Fetching new arrival products")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	products, totalItems, err := h.service.GetNewArrivals(page, limit)
	if err != nil {
		utils.Logger.Printf("❌ [GetNewArrivals] Failed to fetch: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return c.JSON(fiber.Map{
		"current_page": page,
		"items":        products,
		"per_page":     limit,
		"total_items":  totalItems,
		"total_pages":  totalPages,
	})
}
