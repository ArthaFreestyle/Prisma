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

func (s *LecturerServiceImpl) FindAll(c *fiber.Ctx) error {
	ctx := c.UserContext()
	lecturers, err := s.repoLecturer.FindAll(ctx)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}

	response := model.WebResponse[[]model.UserProfile]{
		Status: "success",
		Data:   lecturers,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *LecturerServiceImpl) FindAdvices(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id := c.Params("id")
	lecturer, err := s.repoLecturer.FindAllAdvices(ctx, id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}

	response := model.WebResponse[[]model.UserProfile]{
		Status: "success",
		Data:   lecturer,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
