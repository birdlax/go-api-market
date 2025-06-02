package service

import (
	"backend/domain"
	"math"
	"time"
)

type reportServiceImpl struct {
	repo domain.ReportRepository
}

func NewReportService(repo domain.ReportRepository) domain.ReportService {
	return &reportServiceImpl{repo: repo}
}

func (s *reportServiceImpl) GetRevenueReport(year int, month *int) (map[string]interface{}, error) {
	// --- ดึงข้อมูลช่วงปัจจุบัน ---
	results, totalRevenue, totalOrders, err := s.repo.GetRevenueReport(year, month)
	if err != nil {
		return nil, err
	}

	// คำนวณเปอร์เซ็นต์ของแต่ละหมวดหมู่
	for i, r := range results {
		percentage := (r.TotalRevenue / totalRevenue) * 100
		results[i].PercentageOfTotal = math.Round(percentage*100) / 100
	}

	// --- เตรียมข้อมูลเปรียบเทียบช่วงก่อนหน้า ---
	var prevYear, prevMonth *int
	comparisonType := ""

	if month != nil {
		prev := *month - 1
		if prev == 0 {
			prevMonth = intPtr(12)
			y := year - 1
			prevYear = &y
		} else {
			prevMonth = &prev
			prevYear = &year
		}
		comparisonType = "vs_previous_month"
	} else {
		y := year - 1
		prevYear = &y
		comparisonType = "vs_previous_year"
	}

	_, prevRevenue, prevOrders, err := s.repo.GetRevenueReport(*prevYear, prevMonth)
	if err != nil {
		return nil, err
	}

	report := map[string]interface{}{
		"year":                  year,
		"month":                 month,
		"overall_total_revenue": totalRevenue,
		"overall_total_orders":  totalOrders,
		"report_generated_at":   time.Now().Format(time.RFC3339),
		"revenue_by_category":   results,
		"comparison_type":       comparisonType,
		"previous_period_data": map[string]interface{}{
			"year":                  *prevYear,
			"month":                 prevMonth,
			"overall_total_revenue": prevRevenue,
			"overall_total_orders":  prevOrders,
		},
	}

	return report, nil
}

func intPtr(i int) *int {
	return &i
}

func (s *reportServiceImpl) GetDashboardSummary() (map[string]interface{}, error) {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	startOfWeek := today.AddDate(0, 0, -int(today.Weekday()))
	startOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	todayRevenue, todayOrders, _ := s.repo.GetRevenueAndOrdersBetween(startOfDay, today)
	weekRevenue, weekOrders, _ := s.repo.GetRevenueAndOrdersBetween(startOfWeek, today)
	monthRevenue, monthOrders, _ := s.repo.GetRevenueAndOrdersBetween(startOfMonth, today)

	totalProducts, _ := s.repo.CountTotalProducts()
	newCustomers, _ := s.repo.CountNewCustomersThisMonth()
	pendingOrders, _ := s.repo.CountPendingOrders()
	lowStockItems, _ := s.repo.CountLowStockItems()

	return map[string]interface{}{
		"today_revenue":            todayRevenue,
		"today_orders":             todayOrders,
		"week_revenue":             weekRevenue,
		"week_orders":              weekOrders,
		"month_revenue":            monthRevenue,
		"month_orders":             monthOrders,
		"total_products":           totalProducts,
		"new_customers_this_month": newCustomers,
		"pending_orders_count":     pendingOrders,
		"low_stock_items_count":    lowStockItems,
	}, nil
}

// services/report_service.go
func (s *reportServiceImpl) GetSalesTrend(period string, days int) ([]map[string]interface{}, error) {
	return s.repo.GetSalesTrend(period, days)
}
