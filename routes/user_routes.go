package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
	app.Post("/logout", userHandler.Logout)
	app.Get("/user/:id", userHandler.GetByID)

	admin := app.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)
	admin.Get("/me", userHandler.GetCurrentUser)
	admin.Put("/me/updateprofile", userHandler.UpdateProfile)
	admin.Post("/me/updatepassword", userHandler.UpdatePassword)
	admin.Get("/getall", userHandler.GetAll)
	admin.Delete("/user/:id", userHandler.Delete)
	admin.Put("/user/:id", userHandler.Update)

	users := app.Group("/profile", middleware.JWTMiddleware)
	users.Get("/me", userHandler.GetCurrentUser)
	users.Put("/me/updateprofile", userHandler.UpdateProfile)
	users.Post("/me/updatepassword", userHandler.UpdatePassword)

}
