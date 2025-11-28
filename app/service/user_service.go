package service

import (
	"prisma/app/repository"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	FindById(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func (u UserServiceImpl) Create(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserServiceImpl) Update(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserServiceImpl) Delete(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserServiceImpl) FindById(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserServiceImpl) FindAll(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserServiceImpl) Logout(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserServiceImpl) Login(c *fiber.Ctx) error {
	panic("implement me")
}
