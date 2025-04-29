package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App, productHandler *handler.ProductHandler) {
	app.Use(middleware.CORSMiddleware())
	app.Post("/product", productHandler.CreateProduct)
	app.Get("/product", productHandler.GetAllProduct)
	app.Get("/products", productHandler.GetAllProducts)
	app.Get("/product/:id", productHandler.GetProductByID)
	app.Put("/product/:id", productHandler.UpdateProduct)
	app.Get("/product/:name", productHandler.GetproductByName)
	app.Delete("/product/:id", productHandler.Delete)

	app.Post("/category", productHandler.CreateCategory)
	app.Get("/product/category/:category", productHandler.GetproductByCategory)
}
