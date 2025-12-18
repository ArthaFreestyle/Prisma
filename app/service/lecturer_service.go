package service

import (
	"prisma/app/model"
	"prisma/app/repository"

	"github.com/gofiber/fiber/v2"
)

type LecturerService interface {
	FindByID(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	FindAdvices(c *fiber.Ctx) error
}

type LecturerServiceImpl struct {
	repoLecturer repository.LecturerRepository
	repoStudent  repository.StudentRepository
}

func NewLecturerService(repoLecturer repository.LecturerRepository, repoStudent repository.StudentRepository) LecturerService {
	return &LecturerServiceImpl{
		repoLecturer: repoLecturer,
		repoStudent:  repoStudent,
	}
}

// FindByID godoc
// @Summary Get lecturer by ID
// @Description Retrieve lecturer profile details by ID
// @Tags Lecturers
// @Accept json
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} model.SwaggerWebResponseUserProfile "Successfully retrieved lecturer"
// @Failure 400 {object} model.SwaggerWebResponseString "Lecturer not found"
// @Security BearerAuth
// @Router /lecturers/{id} [get]
func (s *LecturerServiceImpl) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()

	lecturer, err := s.repoLecturer.FindById(ctx, id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := model.WebResponse[*model.UserProfile]{
		Status: "success",
		Data:   lecturer,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// FindAll godoc
// @Summary Get all lecturers
// @Description Retrieve a list of all lecturers
// @Tags Lecturers
// @Accept json
// @Produce json
// @Success 200 {object} model.SwaggerWebResponseUserProfiles "Successfully retrieved lecturers"
// @Failure 400 {object} model.SwaggerWebResponseString "Bad request"
// @Security BearerAuth
// @Router /lecturers [get]
func (s *LecturerServiceImpl) FindAll(c *fiber.Ctx) error {
	ctx := c.UserContext()
	lecturers, err := s.repoLecturer.FindAll(ctx)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := model.WebResponse[[]model.UserProfile]{
		Status: "success",
		Data:   lecturers,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// FindAdvices godoc
// @Summary Get advised students
// @Description Get list of students advised by a specific lecturer
// @Tags Lecturers
// @Accept json
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} model.SwaggerWebResponseUserProfiles "Successfully retrieved advised students"
// @Failure 400 {object} model.SwaggerWebResponseString "Bad request"
// @Security BearerAuth
// @Router /lecturers/{id}/advices [get]
func (s *LecturerServiceImpl) FindAdvices(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	lecturer, err := s.repoLecturer.FindAllAdvices(ctx, id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := model.WebResponse[[]model.UserProfile]{
		Status: "success",
		Data:   lecturer,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
