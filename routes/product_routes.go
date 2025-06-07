package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App, productHandler *handler.ProductHandler) {
	app.Use(middleware.CORSMiddleware())
	api := app.Group("/api")
	api.Get("/products", productHandler.GetAllProducts)
	api.Get("/product/:id", productHandler.GetProductByID)
	api.Get("/categories", productHandler.GetAllCategories)
	api.Get("/product", productHandler.GetAllProduct)
	api.Get("/filter/category/:category", productHandler.GetProductByCategory)
	api.Get("/products/new-arrivals", productHandler.GetNewArrivals)
	api.Get("/categories/filter/:category", productHandler.GetProductByCategory)

	//admin
	admin := api.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)
	admin.Post("/product/bulk", productHandler.CreateMultipleProducts)
	admin.Post("/product/bulkpro", productHandler.CreateMultipleProductsPro)

	admin.Get("/products", productHandler.GetAllProducts)
	admin.Put("/product/:id", productHandler.UpdateProduct)
	admin.Get("/product/:id", productHandler.GetProductByID)
	admin.Get("/product/:name", productHandler.GetProductByName)
	admin.Delete("/product/:id", productHandler.Delete)

	categoryGroup := api.Group("/admin/categories", middleware.JWTMiddleware, middleware.AdminOnly)
	categoryGroup.Get("/", productHandler.GetAllCategories)
	categoryGroup.Get("/filter/:category", productHandler.GetProductByCategory)
	categoryGroup.Post("/bulk", productHandler.CreateMultipleCategories)

}
