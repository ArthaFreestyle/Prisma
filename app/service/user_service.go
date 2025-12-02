package service

import (
	"database/sql"
	"prisma/app/model"
	"prisma/app/repository"
	"prisma/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	FindById(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
}

func NewUserService(repoUser repository.UserRepository, repoStudent repository.StudentRepository, repoLecturer repository.LecturerRepository, DB *sql.DB, validate *validator.Validate, log *logrus.Logger) UserService {
	return &UserServiceImpl{repoUser, repoStudent, repoLecturer, DB, validate, log}
}

type UserServiceImpl struct {
	repoUser     repository.UserRepository
	repoStudent  repository.StudentRepository
	repoLecturer repository.LecturerRepository
	DB           *sql.DB
	validate     *validator.Validate
	Log          *logrus.Logger
}

func (s *UserServiceImpl) Create(c *fiber.Ctx) error {
	var request model.UserCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return fiber.ErrBadRequest
	}

	err := s.validate.Struct(request)
	if err != nil {
		return fiber.ErrBadRequest
	}
	ctx := c.UserContext()
	tx, err := s.DB.Begin()
	PasswordHash, err := utils.HashPassword(request.Password)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	user := &model.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: PasswordHash,
		FullName:     request.FullName,
		RoleId:       request.RoleID,
	}
	user, err = s.repoUser.Save(ctx, user)
	if err != nil {
		err = tx.Rollback()
		return fiber.ErrInternalServerError
	}
	if request.RoleID == "11111111-1111-1111-1111-111111111111" {
		if request.StudentProfile == nil {
			return fiber.ErrBadRequest
		}
		student := model.Student{
			UserID:       user.ID,
			StudentID:    request.StudentProfile.StudentID,
			ProgramStudy: request.StudentProfile.ProgramStudy,
			AcademicYear: request.StudentProfile.AcademicYear,
			AdvisorID:    request.StudentProfile.AdvisorID,
		}
		student, err = s.repoStudent.Save(ctx, tx, student)
		if err != nil {
			err = tx.Rollback()
			return fiber.ErrInternalServerError
		}
	} else if request.RoleID == "22222222-2222-2222-2222-222222222222" {
		if request.LecturerProfile == nil {
			return fiber.ErrBadRequest
		}
		lecturer := &model.Lecturer{
			UserID:     user.ID,
			LecturerID: request.LecturerProfile.LecturerID,
			Department: request.LecturerProfile.Department,
		}
		lecturer, err := s.repoLecturer.Save(ctx, tx, lecturer)
		if err != nil {
			err = tx.Rollback()
			return fiber.ErrInternalServerError
		}
	}

	panic("implement me")
}

func (s *UserServiceImpl) Update(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *UserServiceImpl) Delete(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *UserServiceImpl) FindById(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *UserServiceImpl) FindAll(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
