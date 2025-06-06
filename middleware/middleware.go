package middleware

import (
	"backend/config"
	"backend/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

// func JWTMiddleware(c *fiber.Ctx) error {
// 	// Log Header "Cookie" ดิบๆ ที่ Server ได้รับทันที
// 	rawCookieHeader := c.Get("Cookie")
// 	log.Printf("JWTMiddleware: Raw Cookie Header received: '%s'", rawCookieHeader)

// 	// ดึงค่า Cookie "JWT"
// 	token := c.Cookies(config.JwtCookieName) // ใช้ config.JwtCookieName เพื่อความสอดคล้อง
// 	log.Printf("JWTMiddleware: Value from c.Cookies(\"%s\"): '%s'", config.JwtCookieName, token)

// 	if token == "" {
// 		log.Println("JWTMiddleware: Cookie 'JWT' is empty or not found. Responding 401 Unauthorized.")
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
// 	}

// 	log.Printf("JWTMiddleware: Attempting to parse token: '%s'", token)
// 	claims, err := utils.ParseToken(token) // สมมติว่า utils.ParseToken ใช้ config.JwtSecret
// 	if err != nil {
// 		log.Printf("JWTMiddleware: Error parsing token: %v. Responding 401 Invalid token.", err)
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
// 	}
// 	log.Printf("JWTMiddleware: Token parsed successfully. Claims: %+v", claims)

// 	// ดึง user_id
// 	userIDClaim, userIDOk := claims["user_id"]
// 	if !userIDOk {
// 		log.Printf("JWTMiddleware: 'user_id' claim missing. Claims: %+v", claims)
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token structure: user_id missing"})
// 	}
// 	idFloat, idFloatOk := userIDClaim.(float64)
// 	if !idFloatOk {
// 		log.Printf("JWTMiddleware: 'user_id' claim is not a float64. Type: %T, Value: %v. Claims: %+v", userIDClaim, userIDClaim, claims)
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid ID format in token"})
// 	}
// 	c.Locals("user_id", uint(idFloat))
// 	log.Printf("JWTMiddleware: user_id set in locals: %d", uint(idFloat))

// 	// ดึง role
// 	roleClaim, roleClaimOk := claims["role"]
// 	if !roleClaimOk {
// 		log.Printf("JWTMiddleware: 'role' claim missing. Claims: %+v", claims)
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token structure: role missing"})
// 	}
// 	roleString, roleStringOk := roleClaim.(string)
// 	if !roleStringOk {
// 		log.Printf("JWTMiddleware: 'role' claim is not a string. Type: %T, Value: %v. Claims: %+v", roleClaim, roleClaim, claims)
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role format in token"})
// 	}
// 	c.Locals("role", roleString)
// 	log.Printf("JWTMiddleware: role set in locals: '%s'", roleString)

// 	return c.Next()
// }

func JWTMiddleware(c *fiber.Ctx) error {
	var tokenString string // สร้างตัวแปรว่างสำหรับเก็บ Token

	// --- **2. แก้ไขส่วนนี้: เพิ่ม Logic การอ่านจาก Header** ---
	// มองหา Token จาก "Authorization" Header ก่อน
	authHeader := c.Get("Authorization")
	log.Printf("JWTMiddleware: Authorization Header received: '%s'", authHeader)
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		log.Println("JWTMiddleware: Token found in Authorization Header.")
	}

	// ถ้าใน Header ไม่มี Token ให้ไปมองหาใน Cookie (เป็น Fallback)
	if tokenString == "" {
		rawCookieHeader := c.Get("Cookie")
		log.Printf("JWTMiddleware: Raw Cookie Header received: '%s'", rawCookieHeader)
		tokenString = c.Cookies(config.JwtCookieName)
		if tokenString != "" {
			log.Println("JWTMiddleware: Token found in 'JWT' Cookie.")
		}
	}
	// --------------------------------------------------------

	// --- **3. แก้ไขส่วนนี้: ใช้ตัวแปร tokenString ในการตรวจสอบ** ---
	if tokenString == "" {
		log.Println("JWTMiddleware: No token found in Authorization header or 'JWT' cookie. Responding 401 Unauthorized.")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	// --------------------------------------------------------

	// --- **4. ส่วนที่เหลือ: ใช้ตัวแปร tokenString แทน token เดิม** ---
	log.Printf("JWTMiddleware: Attempting to parse token: '%s'", tokenString)
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		log.Printf("JWTMiddleware: Error parsing token: %v. Responding 401 Invalid token.", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
	log.Printf("JWTMiddleware: Token parsed successfully. Claims: %+v", claims)

	// ส่วนที่เหลือของโค้ดในการดึง user_id และ role เหมือนเดิม...
	userIDClaim, userIDOk := claims["user_id"]
	if !userIDOk {
		log.Printf("JWTMiddleware: 'user_id' claim missing. Claims: %+v", claims)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token structure: user_id missing"})
	}
	idFloat, idFloatOk := userIDClaim.(float64)
	if !idFloatOk {
		log.Printf("JWTMiddleware: 'user_id' claim is not a float64. Type: %T, Value: %v. Claims: %+v", userIDClaim, userIDClaim, claims)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid ID format in token"})
	}
	c.Locals("user_id", uint(idFloat))
	log.Printf("JWTMiddleware: user_id set in locals: %d", uint(idFloat))

	roleClaim, roleClaimOk := claims["role"]
	if !roleClaimOk {
		log.Printf("JWTMiddleware: 'role' claim missing. Claims: %+v", claims)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token structure: role missing"})
	}
	roleString, roleStringOk := roleClaim.(string)
	if !roleStringOk {
		log.Printf("JWTMiddleware: 'role' claim is not a string. Type: %T, Value: %v. Claims: %+v", roleClaim, roleClaim, claims)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role format in token"})
	}
	c.Locals("role", roleString)
	log.Printf("JWTMiddleware: role set in locals: '%s'", roleString)

	return c.Next()
}

func AdminOnly(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin access only"})
	}
	return c.Next()
}
