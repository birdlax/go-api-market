package repository

import (
	"backend/domain"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(product domain.Product) error {
	return r.db.Create(&product).Error
}

func (r *productRepositoryImpl) GetAllProduct() ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.Preload("Images").Order("id ASC").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepositoryImpl) GetProductByName(name string) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Where("name = ?", name).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&domain.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}

// category repository
func (r *productRepositoryImpl) CreateCategory(category domain.Category) error {
	if err := r.db.Create(&category).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepositoryImpl) GetProductByCategory(
	category string,
	page, limit int,
	sort, order string,
	minPrice, maxPrice float64,
) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalItems int64

	offset := (page - 1) * limit

	// Map สำหรับ sort field ที่อนุญาต
	validSortFields := map[string]string{
		"id":        "id",
		"name":      "name",
		"price":     "price",
		"createdat": "created_at",
		"updatedat": "updated_at",
	}

	sortField := validSortFields[strings.ToLower(sort)]
	if sortField == "" {
		sortField = "created_at"
	}

	sortOrder := "ASC"
	if strings.ToLower(order) == "desc" {
		sortOrder = "DESC"
	}

	query := r.db.Model(&domain.Product{}).Where("category_id = ?", category)

	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// นับจำนวนสินค้าทั้งหมด
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// ดึงสินค้าตามเงื่อนไข
	if err := query.
		Preload("Images").
		Preload("Category").
		Order(fmt.Sprintf("%s %s", sortField, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (r *productRepositoryImpl) GetProductByNameAndCategoryID(name string, categoryID uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("name = ? AND category_id = ?", name, categoryID).First(&product).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}
func (r *productRepositoryImpl) GetAllProducts(
	page, limit int,
	sort, order string,
	minPrice, maxPrice float64,
	search string, // เพิ่ม
) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalItems int64
	offset := (page - 1) * limit

	validSortFields := map[string]string{
		"id":        "id",
		"name":      "name",
		"createdat": "created_at",
		"updatedat": "updated_at",
		"price":     "price",
	}
	sortField := validSortFields[strings.ToLower(sort)]
	if sortField == "" {
		sortField = "created_at"
	}
	sortOrder := "ASC"
	if strings.ToLower(order) == "desc" {
		sortOrder = "DESC"
	}

	query := r.db.Model(&domain.Product{})

	// 🔎 เพิ่มการค้นหา
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 🔢 กรองราคา
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// 👇 นับและดึงข้อมูล
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}
	if err := query.
		Preload("Images").
		Preload("Category").
		Order(fmt.Sprintf("%s %s", sortField, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (r *productRepositoryImpl) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *productRepositoryImpl) GetNewArrivals(page, limit int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalItems int64

	offset := (page - 1) * limit

	if err := r.db.Model(&domain.Product{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Preload("Images").
		Preload("Category").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (r *productRepositoryImpl) CreateBulkProducts(products []*domain.Product) error {
	for _, p := range products {
		if err := r.db.Create(&p).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *productRepositoryImpl) GetProductByNameAndCategoryIDPro(name string, categoryID uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("name = ? AND category_id = ?", name, categoryID).First(&product).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &product, err
}

func (r *productRepositoryImpl) CreateBulkProductsPro(products []*domain.Product) error {
	for _, p := range products {
		if err := r.db.Create(&p).Error; err != nil {
			return err
		}
	}
	return nil
}
func (r *productRepositoryImpl) GetProductByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Preload("Images").Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) UpdateProduct(input *domain.UpdateProductInput) (*domain.Product, error) {
	var product domain.Product

	// หา product เดิมพร้อม preload ภาพและหมวดหมู่
	if err := r.db.Preload("Images").Preload("Category").First(&product, input.ID).Error; err != nil {
		return nil, err
	}

	// อัปเดตข้อมูลพื้นฐาน
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Quantity = input.Quantity
	product.CategoryID = input.CategoryID

	// ---- ลบเฉพาะภาพที่ไม่มีใน KeepImagePaths ----
	if len(input.KeepImagePaths) > 0 {
		// ลบรูปทั้งหมดที่ไม่อยู่ใน KeepImagePaths
		if err := r.db.Where("product_id = ? AND path NOT IN ?", product.ID, input.KeepImagePaths).
			Delete(&domain.ProductImage{}).Error; err != nil {
			return nil, err
		}
	}

	// ---- เพิ่มรูปใหม่ ----
	for _, img := range input.Images {
		img.ProductID = product.ID
		if err := r.db.Create(&img).Error; err != nil {
			return nil, err
		}
	}

	// บันทึก
	if err := r.db.Save(&product).Error; err != nil {
		return nil, err
	}

	// โหลดใหม่
	if err := r.db.Preload("Images").Preload("Category").First(&product, product.ID).Error; err != nil {
		return nil, err
	}

	return &product, nil

}

func (r *productRepositoryImpl) FindProductByID(id uint, product *domain.Product) error {
	return r.db.Preload("Images").First(product, id).Error
}

func (r *productRepositoryImpl) DeleteProductImageByID(imageID uint) error {
	return r.db.Delete(&domain.ProductImage{}, imageID).Error
}
