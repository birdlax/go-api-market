package routes

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	app.Use(middleware.CORSMiddleware())
	api := app.Group("/api")
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)
	api.Post("/logout", userHandler.Logout)
	api.Get("/user/:id", userHandler.GetByID)
	api.Get("/gethello", userHandler.GetHello)
	api.Post("/forgot-password", userHandler.ForgotPassword)
	api.Post("/reset-password", userHandler.ResetPassword)

	admin := api.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)
	admin.Get("/me", userHandler.GetCurrentUser)
	admin.Put("/me/updateprofile", userHandler.UpdateProfile)
	admin.Post("/me/updatepassword", userHandler.UpdatePassword)
	admin.Get("/getall", userHandler.GetAll)
	admin.Delete("/user/:id", userHandler.Delete)
	admin.Put("/user/:id", userHandler.UpdateProfilebyId)
	admin.Get("/user/:id", userHandler.GetByID)

	users := app.Group("/profile", middleware.JWTMiddleware)
	users.Get("/me", userHandler.GetCurrentUser)
	users.Put("/me/updateprofile", userHandler.UpdateProfile)
	users.Post("/me/updatepassword", userHandler.UpdatePassword)
}
