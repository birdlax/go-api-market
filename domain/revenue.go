// domain/revenue.go
package domain

import (
	"time"
)

type RevenueResult struct {
	CategoryID        uint    `json:"category_id"`
	CategoryName      string  `json:"category_name"`
	TotalRevenue      float64 `json:"total_revenue"`
	PercentageOfTotal float64 `json:"percentage_of_total"`
}
type ReportRepository interface {
	GetRevenueReport(year int, month *int) ([]RevenueResult, float64, int, error)
	GetRevenueAndOrdersBetween(start, end time.Time) (float64, int, error)
	CountTotalProducts() (int, error)
	CountNewCustomersThisMonth() (int, error)
	CountPendingOrders() (int, error)
	CountLowStockItems() (int, error)
	GetSalesTrend(period string, days int) ([]map[string]interface{}, error)
}
type ReportService interface {
	GetRevenueReport(year int, month *int) (map[string]interface{}, error)
	GetDashboardSummary() (map[string]interface{}, error)
	GetSalesTrend(period string, days int) ([]map[string]interface{}, error)
}
