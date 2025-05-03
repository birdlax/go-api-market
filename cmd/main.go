package main

import (
	"backend/config"
	"backend/domain"
	"backend/handler"
	"backend/repository"
	"backend/routes"
	"backend/service"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	config.ConnectDatabase()
	config.DB.AutoMigrate(&domain.User{}, &domain.Product{}, &domain.Order{}, &domain.OrderItem{})

	utils.StartUserCountLogger(config.DB)

	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	orderRepo := repository.NewOrderRepository(config.DB)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	routes.UserRoutes(app, userHandler)
	routes.ProductRoutes(app, productHandler)
	routes.OrderRoutes(app, orderHandler)
	app.Listen(":3000")
}
