package routes

import (
	"backend/handler"
	// "backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App, productHandler *handler.ProductHandler) {

	app.Post("/product", productHandler.Create)
	app.Delete("/product/:id", productHandler.Delete)
	app.Get("/product", productHandler.GetAllProduct)
	app.Put("/product/:id", productHandler.UpdateProduct)
	app.Get("/product/:name", productHandler.GetproductByName)
}
