package service

import (
	"prisma/app/model"
	"prisma/app/repository"

	"github.com/gofiber/fiber/v2"
)

type AnalyticsService interface {
	Analytics(c *fiber.Ctx) error
	Report(c *fiber.Ctx) error
}

type AnalyticsServiceImpl struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) *AnalyticsServiceImpl {
	return &AnalyticsServiceImpl{repo: repo}
}

func (s *AnalyticsServiceImpl) Analytics(c *fiber.Ctx) error {
	ctx := c.UserContext()
	Analytics, err := s.repo.Statistics(ctx)
	if err != nil {
		response := model.WebResponse[model.Statistics]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := model.WebResponse[*model.Statistics]{
		Status: "success",
		Data:   Analytics,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (s AnalyticsServiceImpl) Report(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	Report, err := s.repo.Reporting(ctx, id)
	if err != nil {
		response := model.WebResponse[model.Statistics]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := model.WebResponse[*model.Statistics]{
		Status: "success",
		Data:   Report,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
