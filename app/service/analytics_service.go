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

// Analytics godoc
// @Summary      Get General Analytics
// @Description  Retrieve general system statistics and analytics data.
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.WebResponse[model.Statistics]
// @Failure      500  {object}  model.WebResponse[model.Statistics]
// @Security     BearerAuth
// @Router       /analytics [get]
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

	response := model.WebResponse[[]model.Statistics]{
		Status: "success",
		Data:   Analytics,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// Report godoc
// @Summary      Get Specific Report
// @Description  Retrieve specific analytics report by ID.
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Report ID"
// @Success      200  {object}  model.WebResponse[model.Statistics]
// @Failure      500  {object}  model.WebResponse[model.Statistics]
// @Security     BearerAuth
// @Router       /analytics/report/{id} [get]
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

	response := model.WebResponse[[]*model.Statistics]{
		Status: "success",
		Data:   Report,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
