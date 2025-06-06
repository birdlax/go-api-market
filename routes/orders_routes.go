package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(app *fiber.App, orderHandler *handler.OrderHandler) {
	app.Use(middleware.CORSMiddleware())
	user := app.Group("/order", middleware.JWTMiddleware)
	user.Get("/:id", orderHandler.GetOrderByID)
	user.Get("/", orderHandler.GetOrder)
	user.Get("/show/orderalls", orderHandler.GetOrdersByStatus)
	user.Put("/pay", orderHandler.MarkOrderAsPaid)
	user.Put("/cancel", orderHandler.CancelOrder)

	//admin
	admin := app.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)
	admin.Get("/orders", orderHandler.GetAllOrders)
	admin.Get("/order/:id", orderHandler.GetOrderByID)
	admin.Put("/order/:id", orderHandler.UpdateOrder)
	admin.Delete("/order/:id", orderHandler.DeleteOrder)

}
