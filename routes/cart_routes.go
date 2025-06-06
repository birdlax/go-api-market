package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func CartRoutes(app *fiber.App, cartHandler *handler.CartHandler) {
	cart := app.Group("/cart", middleware.JWTMiddleware)
	cart.Post("/item", cartHandler.AddItem)
	cart.Delete("/items/:product_id", cartHandler.RemoveItem)
	cart.Get("/", cartHandler.GetCart) // <--- แก้จาก "/cart" เป็น "/"
	cart.Post("/checkout", cartHandler.Checkout)
	cart.Delete("/itemx/:product_id", cartHandler.RemoveItemOne)
	cart.Post("/item/:product_id", cartHandler.AddOneItem)
}
