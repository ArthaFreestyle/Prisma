package routes

import (
	"prisma/app/service"
	"prisma/middleware"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	AuthService        service.AuthService
	UserService        service.UserService
	AchievementService service.AchievementService
	StudentService     service.StudentService
	LecturerService    service.LecturerService
	AnalyticsService   service.AnalyticsService
	AuthMiddleware     fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/v1/auth/login", c.AuthService.Login)
	c.App.Post("/api/v1/auth/refresh", c.AuthService.RefreshToken)
	c.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	c.App.Post("/api/v1/auth/logout", c.AuthService.Logout)
	c.App.Get("/api/v1/auth/profile", c.UserService.Profile)

	//users
	c.App.Post("/api/v1/users", middleware.RequirePermission("users:create"), c.UserService.Create)
	c.App.Get("/api/v1/users", middleware.RequirePermission("users:list"), c.UserService.FindAll)
	c.App.Get("/api/v1/users/:id", middleware.RequirePermission("users:detail"), c.UserService.FindById)
	c.App.Put("/api/v1/users/:id", middleware.RequirePermission("users:update"), c.UserService.Update)
	c.App.Delete("/api/v1/users/:id", middleware.RequirePermission("users:delete"), c.UserService.Delete)
	c.App.Put("/api/v1/users/:id/role", middleware.RequirePermission("users:updateRole"), c.UserService.UpdateRole)

	//achievement
	c.App.Post("/api/v1/achievements", middleware.RequirePermission("achievements:create"), c.AchievementService.Create)
	c.App.Get("/api/v1/achievements", middleware.RequirePermission("achievements:list"), c.AchievementService.FindAll)
	c.App.Get("/api/v1/achievements/:id", middleware.RequirePermission("achievements:detail"), c.AchievementService.FindByID)
	c.App.Put("/api/v1/achievements/:id", middleware.RequirePermission("achievements:update"), c.AchievementService.Update)
	c.App.Delete("/api/v1/achievements/:id", middleware.RequirePermission("achievements:delete"), c.AchievementService.Delete)
	c.App.Post("/api/v1/achievements:id/submit", middleware.RequirePermission("achievements:submit"), c.AchievementService.Submit)
	c.App.Post("/api/achievements/:id/verify", middleware.RequirePermission("achievements:verify"), c.AchievementService.Verify)
	c.App.Post("/api/v1/achievements/:id/reject", middleware.RequirePermission("achievements:reject"), c.AchievementService.Reject)
	c.App.Get("/api/v1/achievements/:id/history", middleware.RequirePermission("achievements:history"), c.AchievementService.History)
	c.App.Post("/api/v1/achievements/:id/attachment", middleware.RequirePermission("achievements:upload"), c.AchievementService.Attachment)

	//Student And Lecturer
	c.App.Get("/api/v1/students", middleware.RequirePermission("students:list"), c.StudentService.FindAll)
	c.App.Get("/api/v1/students/:id", middleware.RequirePermission("students:detail"), c.StudentService.FindById)
	c.App.Get("/api/v1/students:id/achievements", middleware.RequirePermission("students:achievements"), c.StudentService.FindAchievements)
	c.App.Put("/api/v1/students/:id/advisor", middleware.RequirePermission("students:updateAdvisor"), c.StudentService.ChangeAdvisor)
	c.App.Get("/api/v1/lecturers", middleware.RequirePermission("lecturers:list"), c.LecturerService.FindAll)
	c.App.Get("/api/v1/lecturer/:id", middleware.RequirePermission("lecturers:details"), c.LecturerService.FindByID)
	c.App.Get("/api/v1/lecturers/:id/advices", middleware.RequirePermission("lecturers:advisees"), c.LecturerService.FindAdvices)

	//analytics And Reporting
	c.App.Get("/api/v1/reports/statistics", middleware.RequirePermission("reports:statistics"), c.AnalyticsService.Analytics)
	c.App.Get("/api/v1/reports/student/:id", middleware.RequirePermission("reports:studentDetail"), c.AnalyticsService.Report)
}
