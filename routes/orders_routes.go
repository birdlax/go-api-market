package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(app *fiber.App, orderHandler *handler.OrderHandler) {
	app.Use(middleware.CORSMiddleware())

	app.Use(middleware.JWTMiddleware)
	//admin
	app.Get("/orders", orderHandler.GetAllOrders)
	app.Get("/order/:id", orderHandler.GetOrderByID)
	app.Put("/order/:id", orderHandler.UpdateOrder)
	app.Delete("/order/:id", orderHandler.DeleteOrder)

	app.Get("/order", orderHandler.GetOrder)
	app.Get("/orderalls", orderHandler.GetOrdersByStatus)
	app.Put("/pay", orderHandler.MarkOrderAsPaid)
	app.Put("/orders/cancel", orderHandler.CancelOrder)

	app.Get("/revenue", orderHandler.GetRevenueByCategory)

}
