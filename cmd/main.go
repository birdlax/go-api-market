package main

import (
	"backend/config"
	"backend/domain"
	"backend/handler"
	"backend/repository"
	"backend/routes"
	"backend/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	config.ConnectDatabase()
	config.DB.AutoMigrate(&domain.User{}, &domain.Product{})

	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	routes.UserRoutes(app, userHandler)
	routes.ProductRoutes(app, productHandler)

	app.Listen(":3000")
}
