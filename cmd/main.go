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
	config.DB.AutoMigrate(&domain.User{}, &domain.Product{}, &domain.Order{}, &domain.OrderItem{}, &domain.Cart{}, &domain.CartItem{})

	utils.InitLogger()

	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	orderRepo := repository.NewOrderRepository(config.DB)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	cartRepo := repository.NewCartRepository(config.DB)
	cartService := service.NewCartService(cartRepo, orderService)
	cartHandler := handler.NewCartHandler(cartService)

	routes.UserRoutes(app, userHandler)
	routes.ProductRoutes(app, productHandler)
	routes.OrderRoutes(app, orderHandler)
	routes.CartRoutes(app, cartHandler)
	app.Listen(":3000")
}
