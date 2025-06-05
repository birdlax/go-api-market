package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware() fiber.Handler {
	ipv6Origin := "http://[2403:6200:8833:6faa:7486:aa70:a908:5b2e]"
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173," + ipv6Origin,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowCredentials: true,
	})
}
