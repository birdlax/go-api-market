package main

import (
	"backend/config"
	"backend/handler"
	"backend/middleware"
	"backend/repository"
	"backend/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	config.ConnectDatabase()

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
	app.Post("/logout", userHandler.Logout)

	app.Get("/user/:id", userHandler.GetByID)
	app.Post("/updatepassword/:id", userHandler.UpdatePassword)

	admin := app.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)
	admin.Get("/getall", userHandler.GetAll)
	admin.Delete("/user/:id", userHandler.Delete)
	admin.Put("/user/:id", userHandler.Update)

	users := app.Group("/profile", middleware.JWTMiddleware)
	users.Get("/me", userHandler.GetCurrentUser)

	app.Listen(":3000")
}
