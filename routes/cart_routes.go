package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func CartRoutes(app *fiber.App, cartHandler *handler.CartHandler) {
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.JWTMiddleware)
	app.Post("/cart/item", cartHandler.AddItem)
	app.Delete("/cart/items/:product_id", cartHandler.RemoveItem)
	app.Get("/cart", cartHandler.GetCart)
	app.Post("/cart/checkout", cartHandler.Checkout)
	app.Delete("/cart/itemx/:product_id", cartHandler.RemoveItemOne)
	app.Post("/cart/item/:product_id", cartHandler.AddOneItem)
}
