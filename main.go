package main

import (
	"backend/config"
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

	utils.InitLogger()

	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)
	// caet and order
	orderRepo := repository.NewOrderRepository(config.DB)
	cartRepo := repository.NewCartRepository(config.DB)
	orderService := service.NewOrderService(orderRepo, cartRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	cartService := service.NewCartService(cartRepo, orderService)
	cartHandler := handler.NewCartHandler(cartService)

	addressRepo := repository.NewAddressRepository(config.DB)
	addressService := service.NewAddressService(addressRepo)
	addressHandler := handler.NewAddressHandler(addressService)

	routes.UserRoutes(app, userHandler)
	routes.ProductRoutes(app, productHandler)
	routes.OrderRoutes(app, orderHandler)
	routes.CartRoutes(app, cartHandler)
	routes.AddressRouter(app, addressHandler)
	app.Listen(":3000")
}
