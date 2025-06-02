package repository

import (
	"backend/domain"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type reportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) domain.ReportRepository {
	return &reportRepositoryImpl{db: db}
}

func (r *reportRepositoryImpl) GetRevenueReport(year int, month *int) ([]domain.RevenueResult, float64, int, error) {
	var results []domain.RevenueResult
	var totalRevenue float64
	var totalOrders64 int64

	// === 1. Query รายได้ตามหมวดหมู่ ===
	db := r.db.
		Table("order_items").
		Select(`categories.id AS category_id,
		        categories.name AS category_name,
		        SUM(order_items.quantity * order_items.price) AS total_revenue`).
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN categories ON categories.id = products.category_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("EXTRACT(YEAR FROM orders.created_at) = ?", year)

	if month != nil {
		db = db.Where("EXTRACT(MONTH FROM orders.created_at) = ?", *month)
	}

	// รายได้แยกตามหมวดหมู่
	if err := db.Group("categories.id, categories.name").Scan(&results).Error; err != nil {
		return nil, 0, 0, err
	}

	// === 2. Query รายได้รวมทั้งหมด ===
	revenueQuery := r.db.
		Table("order_items").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("EXTRACT(YEAR FROM orders.created_at) = ?", year)

	if month != nil {
		revenueQuery = revenueQuery.Where("EXTRACT(MONTH FROM orders.created_at) = ?", *month)
	}

	var totalRevenueNull sql.NullFloat64
	if err := revenueQuery.Select("SUM(order_items.quantity * order_items.price)").Scan(&totalRevenueNull).Error; err != nil {
		return nil, 0, 0, err
	}
	if totalRevenueNull.Valid {
		totalRevenue = totalRevenueNull.Float64
	} else {
		totalRevenue = 0
	}

	// === 3. Query จำนวนออเดอร์ ===
	orderQuery := r.db.Model(&domain.Order{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", year)

	if month != nil {
		orderQuery = orderQuery.Where("EXTRACT(MONTH FROM created_at) = ?", *month)
	}

	if err := orderQuery.Count(&totalOrders64).Error; err != nil {
		return nil, 0, 0, err
	}

	return results, totalRevenue, int(totalOrders64), nil
}

func (r *reportRepositoryImpl) GetRevenueAndOrdersBetween(start, end time.Time) (float64, int, error) {
	var revenueNull sql.NullFloat64
	var count int64

	revenueQuery := r.db.Table("order_items").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.created_at BETWEEN ? AND ?", start, end).
		Select("SUM(order_items.quantity * order_items.price)")

	if err := revenueQuery.Scan(&revenueNull).Error; err != nil {
		return 0, 0, err
	}

	if err := r.db.Model(&domain.Order{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&count).Error; err != nil {
		return 0, 0, err
	}

	revenue := 0.0
	if revenueNull.Valid {
		revenue = revenueNull.Float64
	}

	return revenue, int(count), nil
}

func (r *reportRepositoryImpl) CountTotalProducts() (int, error) {
	var count int64
	if err := r.db.Model(&domain.Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *reportRepositoryImpl) CountNewCustomersThisMonth() (int, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var count int64
	if err := r.db.Model(&domain.User{}).
		Where("created_at >= ?", startOfMonth).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *reportRepositoryImpl) CountPendingOrders() (int, error) {
	var count int64
	if err := r.db.Model(&domain.Order{}).
		Where("status = ?", "pending").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *reportRepositoryImpl) CountLowStockItems() (int, error) {
	var count int64
	if err := r.db.Model(&domain.Product{}).
		Where("stock <= ?", 5).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// repositories/report_repository.go
func (r *reportRepositoryImpl) GetSalesTrend(period string, days int) ([]map[string]interface{}, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	var results []map[string]interface{}

	rows, err := r.db.Raw(`
		SELECT
			DATE(orders.created_at) AS date,
			SUM(order_items.quantity * order_items.price) AS revenue,
			COUNT(DISTINCT orders.id) AS order_count
		FROM orders
		JOIN order_items ON order_items.order_id = orders.id
		WHERE DATE(orders.created_at) BETWEEN ? AND ?
		GROUP BY DATE(orders.created_at)
		ORDER BY DATE(orders.created_at)
	`, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var revenue float64
		var orderCount int

		if err := rows.Scan(&date, &revenue, &orderCount); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"date":        date,
			"revenue":     revenue,
			"order_count": orderCount,
		})
	}

	return results, nil
}
