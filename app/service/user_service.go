package service

import (
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
	Logout(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}

func NewUserService(repo repository.UserRepository, logout repository.AuthRepository, Log *logrus.Logger, secret []byte) UserService {
	return &UserServiceImpl{
		repo:   repo,
		Log:    Log,
		Auth:   logout,
		secret: secret,
	}
}

type UserServiceImpl struct {
	repo     repository.UserRepository
	Auth     repository.AuthRepository
	validate *validator.Validate
	Log      *logrus.Logger
	secret   []byte
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
	refreshToken := c.Cookies("refresh_token")
	ctx := c.UserContext()
	err := s.Auth.Logout(ctx, refreshToken)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	response := model.LogoutResponse{
		Message: "Logged out",
	}

	return c.JSON(model.WebResponse[model.LogoutResponse]{
		Data:   response,
		Status: "success",
	})

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

	return c.JSON(model.WebResponse[*model.LoginResponse]{
		Data:   response,
		Status: "success",
	})
}

func (s *UserServiceImpl) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	Claims, err := utils.ValidateToken(refreshToken, s.secret)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	ctx := c.UserContext()
	Access, err := s.Auth.RefreshToken(ctx, refreshToken, s.secret)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	AuthResponse := &model.UserAuthResponse{
		ID:          Claims.UserID,
		FullName:    Claims.FullName,
		Username:    Claims.Username,
		Role:        Claims.Role,
		Permissions: Claims.Permissions,
	}

	response := &model.LoginResponse{
		Token:        Access,
		RefreshToken: refreshToken,
		User:         *AuthResponse,
	}

	return c.JSON(model.WebResponse[*model.LoginResponse]{
		Data:   response,
		Status: "success",
	})

}
