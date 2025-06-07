package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(app *fiber.App, orderHandler *handler.OrderHandler) {
	app.Use(middleware.CORSMiddleware())

	app.Use(middleware.JWTMiddleware)
	api := app.Group("/api")
	api.Get("/order/:id", orderHandler.GetOrderByID)
	api.Get("/order", orderHandler.GetOrder)
	api.Get("/orderalls", orderHandler.GetOrdersByStatus)
	api.Put("/pay", orderHandler.MarkOrderAsPaid)
	api.Put("/orders/cancel", orderHandler.CancelOrder)

	//admin
	admin := api.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)
	admin.Get("/orders", orderHandler.GetAllOrders)
	admin.Get("/order/:id", orderHandler.GetOrderByID)
	admin.Put("/order/:id", orderHandler.UpdateOrder)
	admin.Delete("/order/:id", orderHandler.DeleteOrder)

}
