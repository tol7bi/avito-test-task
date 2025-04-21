package middleware

import (
	"os"
	"strings"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"message": "invalid token"})
		}

		c.Locals("user", token)
		return c.Next()
	}
}

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		userRole := fmt.Sprintf("%v", claims["role"]) 


		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "forbidden",
			})
		}
		return c.Next()
	}
}