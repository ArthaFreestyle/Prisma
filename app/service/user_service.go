package service

import (
	"prisma/app/model"
	"prisma/app/repository"

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

func NewUserService(repo repository.UserRepository, validate *validator.Validate, log *logrus.Logger) UserService {
	return &UserServiceImpl{repo, validate, log}
}

type UserServiceImpl struct {
	repo     repository.UserRepository
	validate *validator.Validate
	Log      *logrus.Logger
}

func (s *UserServiceImpl) Create(c *fiber.Ctx) error {
	var request model.CreateUserRequest
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
