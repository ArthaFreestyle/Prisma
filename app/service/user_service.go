package service

import (
	"prisma/app/model"
	"prisma/app/repository"
	"prisma/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
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

func NewUserService(repo repository.UserRepository, Log *logrus.Logger, secret []byte) UserService {
	return &UserServiceImpl{
		repo:   repo,
		Log:    Log,
		secret: secret,
	}
}

type UserServiceImpl struct {
	repo   repository.UserRepository
	Log    *logrus.Logger
	secret []byte
}

func (s *UserServiceImpl) Create(c *fiber.Ctx) error {
	//TODO implement me
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

func (s *UserServiceImpl) Logout(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *UserServiceImpl) Login(c *fiber.Ctx) error {
	var request = new(model.LoginRequest)
	if err := c.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}
	ctx := c.UserContext()
	User, err := s.repo.FindByUsername(ctx, request.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !utils.CheckPasswordHash(request.Password, User.PasswordHash) {
		return fiber.ErrUnauthorized
	}
	access, refresh, err := utils.GenerateToken(User, s.secret)
	if err != nil {

		return fiber.ErrInternalServerError
	}

	AuthResponse := &model.UserAuthResponse{
		ID:          User.ID,
		FullName:    User.FullName,
		Username:    User.Username,
		Role:        User.RoleName,
		Permissions: User.Permissions,
	}

	response := &model.LoginResponse{
		Token:        access,
		RefreshToken: refresh,
		User:         *AuthResponse,
	}

	return c.JSON(model.WebResponse[*model.LoginResponse]{Data: response})
}
