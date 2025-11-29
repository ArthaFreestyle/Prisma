package routes

import (
	"prisma/app/service"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserService    service.UserService
	AuthMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/login", c.UserService.Login)
	c.App.Post("/api/refresh", c.UserService.RefreshToken)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Post("/api/logout", c.UserService.Logout)
}
