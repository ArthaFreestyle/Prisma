package middleware

import (
	"context"
	"prisma/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(JWTsecret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses tidak valid",
			})
		}

		claims, err := utils.ValidateToken(tokenParts[1], JWTsecret)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses tidak valid",
			})
		}

		ctx := context.WithValue(c.UserContext(), "user", claims)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func RequirePermission(requiredPerm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Locals("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)

		rawPerms, ok := claims["permissions"].([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Permissions not found in token",
			})
		}

		hasPermission := false
		for _, p := range rawPerms {
			if strPerm, ok := p.(string); ok {
				if strPerm == requiredPerm {
					hasPermission = true
					break
				}
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "You don't have permission: " + requiredPerm,
			})
		}

		return c.Next()
	}
}
