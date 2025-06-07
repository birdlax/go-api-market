package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func CartRoutes(app *fiber.App, cartHandler *handler.CartHandler) {
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.JWTMiddleware)
	api := app.Group("/api")
	api.Post("/cart/item", cartHandler.AddItem)
	api.Delete("/cart/items/:product_id", cartHandler.RemoveItem)
	api.Get("/cart", cartHandler.GetCart)
	api.Post("/cart/checkout", cartHandler.Checkout)
	api.Delete("/cart/itemx/:product_id", cartHandler.RemoveItemOne)
	api.Post("/cart/item/:product_id", cartHandler.AddOneItem)
}
