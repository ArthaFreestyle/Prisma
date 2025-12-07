package routes

import (
	"prisma/app/service"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	AuthService        service.AuthService
	UserService        service.UserService
	AchievementService service.AchievementService
	AuthMiddleware     fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/v1/auth/login", c.AuthService.Login)
	c.App.Post("/api/v1/auth/refresh", c.AuthService.RefreshToken)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Post("/api/v1/logout", c.AuthService.Logout)
	c.App.Get("/api/v1/auth/profile", c.UserService.Profile)

	//users
	c.App.Post("/api/v1/users", c.UserService.Create)
	c.App.Get("/api/v1/users", c.UserService.FindAll)
	c.App.Get("/api/v1/users/:id", c.UserService.FindById)
	c.App.Put("/api/v1/users/:id", c.UserService.Update)
	c.App.Delete("/api/v1/users/:id", c.UserService.Delete)
	c.App.Put("/api/v1/users/:id/role", c.UserService.UpdateRole)

	//achievement
	c.App.Post("/api/v1/achievements", c.AchievementService.Create)

}
