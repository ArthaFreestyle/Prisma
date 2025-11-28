package routes

import (
	"prisma/app/service"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserService    service.UserService
	AuthMiddleware *fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRole()
}

func (c *RouteConfig) SetupGuestRole() {
	c.App.Post("/api/login", c.UserService.Login)
}
