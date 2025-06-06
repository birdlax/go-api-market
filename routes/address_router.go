package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func AddressRouter(app *fiber.App, addressHandler *handler.AddressHandler) {
	app.Use(middleware.CORSMiddleware())
	addressGroup := app.Group("/api/addresses", middleware.JWTMiddleware)
	addressGroup.Post("/", addressHandler.CreateAddress)
	addressGroup.Get("/", addressHandler.GetAddresses)
	addressGroup.Get("/:id", addressHandler.GetAddressByID)
	addressGroup.Put("/update/:id", addressHandler.UpdateAddress)
	addressGroup.Delete("/delete/:id", addressHandler.DeleteAddress)
	addressGroup.Put("/default/:id", addressHandler.SwitchDefault)

}
