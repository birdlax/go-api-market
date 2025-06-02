package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App, reportHandler *handler.ReportHandler) {
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.JWTMiddleware)
	report := app.Group("/admin/reports", middleware.JWTMiddleware, middleware.AdminOnly)
	report.Get("/revenue", reportHandler.GetRevenueReport)
	report.Get("/dashboard-summary", reportHandler.GetDashboardSummary)
	report.Get("/sales-trend", reportHandler.GetSalesTrend)
}
