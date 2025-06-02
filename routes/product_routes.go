package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App, productHandler *handler.ProductHandler) {
	app.Use(middleware.CORSMiddleware())
	app.Get("/products", productHandler.GetAllProducts)
	app.Get("/product/:id", productHandler.GetProductByID)
	app.Get("/categories", productHandler.GetAllCategories)
	app.Get("/product", productHandler.GetAllProduct)
	app.Get("/filter/category/:category", productHandler.GetProductByCategory)
	app.Get("/products/new-arrivals", productHandler.GetNewArrivals)
	app.Get("/categories/filter/:category", productHandler.GetProductByCategory)
	//admin
	admin := app.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)
	admin.Post("/product/bulk", productHandler.CreateMultipleProducts)

	admin.Get("/products", productHandler.GetAllProducts)
	admin.Put("/product/:id", productHandler.UpdateProduct)
	admin.Get("/product/:id", productHandler.GetProductByID)
	admin.Get("/product/:name", productHandler.GetProductByName)
	admin.Delete("/product/:id", productHandler.Delete)

	categoryGroup := app.Group("/admin/categories", middleware.JWTMiddleware, middleware.AdminOnly)
	categoryGroup.Get("/", productHandler.GetAllCategories)
	categoryGroup.Get("/filter/:category", productHandler.GetProductByCategory)
	categoryGroup.Post("/bulk", productHandler.CreateMultipleCategories)

}
