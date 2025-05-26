package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func AddressRouter(app *fiber.App, addressHandler *handler.AddressHandler) {
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.JWTMiddleware)

	addressGroup := app.Group("/api/addresses")
	addressGroup.Post("/", addressHandler.CreateAddress)
	addressGroup.Get("/", addressHandler.GetAddresses)
	addressGroup.Put("/:id", addressHandler.UpdateAddress)
	addressGroup.Delete("/:id", addressHandler.DeleteAddress)

}
