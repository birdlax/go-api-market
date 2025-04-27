package routes

import (
	"backend/handler"
	"backend/middleware"
	// "backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App, productHandler *handler.ProductHandler) {
	app.Get("/product", productHandler.GetAllProduct)
	app.Use(middleware.CORSMiddleware())
	app.Post("/product", productHandler.CreateProduct)
	app.Put("/product/:id", productHandler.UpdateProduct)
	app.Get("/product/:name", productHandler.GetproductByName)
	app.Delete("/product/:id", productHandler.Delete)
	app.Get("/product/category/:category", productHandler.GetproductByCategory)
	app.Post("/category", productHandler.CreateCategory)
}
