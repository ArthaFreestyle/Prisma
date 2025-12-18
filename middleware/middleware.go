package middleware

import (
	"context"
	"prisma/app/model"
	"prisma/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
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
		// 1. Safety check: Ensure the context value exists
		userVal := c.UserContext().Value("user")
		if userVal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User context not found",
			})
		}

		// 2. Type Assertion: Ensure it's the correct Claims struct
		claims, ok := userVal.(*model.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token structure",
			})
		}

		// 3. Direct Check: No need for interface assertion, just loop the []string
		hasPermission := false
		for _, p := range claims.Permissions {
			if p == requiredPerm {
				hasPermission = true
				break
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
