package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware() fiber.Handler {
	ipv6Origin := "http://223.206.229.144"
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173," + ipv6Origin,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowCredentials: true,
	})
}
