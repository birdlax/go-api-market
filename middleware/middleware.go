package middleware

import (
	"backend/utils"
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("JWT")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	if idFloat, ok := claims["user_id"].(float64); ok {
		c.Locals("user_id", uint(idFloat))
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	c.Locals("role", claims["role"])
	// ใน JWTMiddleware บน VM
	allCookies := c.GetReqHeaders()["Cookie"]
	utils.Logger.Printf("JWTMiddleware: All cookies received by server: %s", allCookies)

	token = c.Cookies("JWT")
	utils.Logger.Printf("🎉 JWTMiddleware: Token from c.Cookies(\"JWT\"): '%s'", token)
	return c.Next()
}

func AdminOnly(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin access only"})
	}
	return c.Next()
}
