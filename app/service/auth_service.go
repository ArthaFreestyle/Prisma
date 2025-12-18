package service

import (
	"prisma/app/model"
	"prisma/app/repository"
	"prisma/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthService interface {
	Logout(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}

func NewAuthService(repo repository.UserRepository, logout repository.AuthRepository, Log *logrus.Logger, secret []byte) AuthService {
	return &AuthServiceImpl{
		repo:   repo,
		Log:    Log,
		Auth:   logout,
		secret: secret,
	}
}

type AuthServiceImpl struct {
	repo     repository.UserRepository
	Auth     repository.AuthRepository
	validate *validator.Validate
	Log      *logrus.Logger
	secret   []byte
}

// Logout godoc
// @Summary      Logout User
// @Description  Invalidate refresh token and logout user.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.WebResponse[model.LogoutResponse]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /auth/logout [post]
func (s *AuthServiceImpl) Logout(c *fiber.Ctx) error {
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

// Login godoc
// @Summary      Login User
// @Description  Authenticate user and return access/refresh tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequest true "Login Credentials"
// @Success      200  {object}  model.WebResponse[model.LoginResponse]
// @Failure      400  {object}  model.WebResponse[string]
// @Failure      401  {object}  model.WebResponse[string]
// @Failure      500  {object}  model.WebResponse[string]
// @Router       /auth/login [post]
func (s *AuthServiceImpl) Login(c *fiber.Ctx) error {
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

// RefreshToken godoc
// @Summary      Refresh Access Token
// @Description  Get a new access token using a valid refresh token from cookies.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.WebResponse[model.LoginResponse]
// @Failure      401  {object}  model.WebResponse[string]
// @Router       /auth/refresh [post]
func (s *AuthServiceImpl) RefreshToken(c *fiber.Ctx) error {
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
