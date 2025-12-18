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

// FindAll godoc
// @Summary Get all students
// @Description Retrieve a list of all students profiles
// @Tags Students
// @Accept json
// @Produce json
// @Success 200 {object} model.SwaggerWebResponseUserProfiles "Successfully retrieved students"
// @Failure 400 {object} model.SwaggerWebResponseString "Bad request"
// @Security BearerAuth
// @Router /students [get]
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

// FindById godoc
// @Summary Get student by ID
// @Description Retrieve detailed student profile by ID
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} model.SwaggerWebResponseUserProfile "Successfully retrieved student"
// @Failure 400 {object} model.SwaggerWebResponseString "Student not found"
// @Security BearerAuth
// @Router /students/{id} [get]
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

// FindAchievements godoc
// @Summary Get student achievements
// @Description Retrieve list of achievements belonging to a student with pagination
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {object} model.SwaggerWebResponseAchievementReferenceAdmin "Successfully retrieved achievements"
// @Failure 400 {object} model.SwaggerWebResponseString "Bad request"
// @Security BearerAuth
// @Router /students/{id}/achievements [get]
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

// ChangeAdvisor godoc
// @Summary Change student advisor
// @Description Assign or change the academic advisor for a student
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body model.SwaggerChangeAdvisorRequest true "Advisor ID payload"
// @Success 200 {object} model.SwaggerWebResponseStudent "Successfully changed advisor"
// @Failure 400 {object} model.SwaggerWebResponseString "Bad request"
// @Security BearerAuth
// @Router /students/{id}/advisor [patch]
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
