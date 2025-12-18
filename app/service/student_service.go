package service

import (
	"database/sql"
	"prisma/app/model"
	"prisma/app/repository"

	"github.com/gofiber/fiber/v2"
)

type StudentService interface {
	FindAll(c *fiber.Ctx) error
	FindById(c *fiber.Ctx) error
	FindAchievements(c *fiber.Ctx) error
	ChangeAdvisor(c *fiber.Ctx) error
}

type StudentServiceImpl struct {
	repoStudent     repository.StudentRepository
	repoAchievement repository.AchievementReferenceRepository
}

func NewStudentService(repoStudent repository.StudentRepository, repoAchievement repository.AchievementReferenceRepository) StudentService {
	return &StudentServiceImpl{
		repoStudent:     repoStudent,
		repoAchievement: repoAchievement,
	}
}

func (s *StudentServiceImpl) FindAll(c *fiber.Ctx) error {
	ctx := c.UserContext()
	Students, err := s.repoStudent.FindAll(ctx)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := model.WebResponse[[]model.UserProfile]{
		Status: "success",
		Data:   Students,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *StudentServiceImpl) FindById(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()
	Student, err := s.repoStudent.FindById(ctx, id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := model.WebResponse[*model.UserProfile]{
		Status: "success",
		Data:   Student,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *StudentServiceImpl) FindAchievements(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	Page := c.QueryInt("page", 1)
	Limit := c.QueryInt("limit", 10)
	Achievements, err := s.repoAchievement.FindByStudentId(ctx, id, Page, Limit)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	response := model.WebResponse[[]model.AchievementReferenceAdmin]{
		Status: "success",
		Data:   Achievements,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *StudentServiceImpl) ChangeAdvisor(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()
	var req struct {
		AdvisorID string `json:"advisor"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data JSON tidak valid",
		})
	}

	Student, err := s.repoStudent.FindById(ctx, id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	if req.AdvisorID != "" {
		Student.AdvisorID = sql.NullString{
			String: req.AdvisorID,
			Valid:  true,
		}
	} else {
		Student.AdvisorID = sql.NullString{
			Valid: false,
		}
	}
	StudentUpdate := &model.Student{
		AdvisorID: req.AdvisorID,
		ID:        id,
	}

	StudentAfterUpdate, err := s.repoStudent.UpdateById(ctx, StudentUpdate)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := model.WebResponse[*model.Student]{
		Status: "success",
		Data:   StudentAfterUpdate,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
