package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(app *fiber.App, orderHandler *handler.OrderHandler) {
	app.Use(middleware.CORSMiddleware())

	app.Get("/order/:id", orderHandler.GetOrderByID)
	app.Get("/order", orderHandler.GetOrder)
	app.Get("/orderalls", orderHandler.GetOrdersByStatus)
	app.Put("/pay", orderHandler.MarkOrderAsPaid)
	app.Put("/orders/cancel", orderHandler.CancelOrder)

	//admin
	admin := app.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)
	admin.Get("/orders", orderHandler.GetAllOrders)
	admin.Get("/order/:id", orderHandler.GetOrderByID)
	admin.Put("/order/:id", orderHandler.UpdateOrder)
	admin.Delete("/order/:id", orderHandler.DeleteOrder)

}
