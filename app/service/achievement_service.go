package service

import (
	"prisma/app/model"
	"prisma/app/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AchievementService interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	FindByID(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
	Submit(c *fiber.Ctx) error
	History(c *fiber.Ctx) error
	Attachment(c *fiber.Ctx) error
}

type AchievementServiceImpl struct {
	repoAchievement repository.AchievementRepository
	validate        *validator.Validate
	Log             *logrus.Logger
}

func NewAchievementService(repo repository.AchievementRepository, validate *validator.Validate, Log *logrus.Logger) *AchievementServiceImpl {
	return &AchievementServiceImpl{
		repoAchievement: repo,
		validate:        validate,
		Log:             Log,
	}
}

func (s *AchievementServiceImpl) Create(c *fiber.Ctx) error {
	var request model.CreateAchievementRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data":   err.Error(),
		})
	}
	ctx := c.UserContext()
	val := ctx.Value("user")
	achievement := &model.AchievementMongo{
		StudentID:       val.(*model.Claims).UserID,
		AchievementType: request.AchievementType,
		Title:           request.Title,
		Description:     request.Description,
		Details:         request.Details,
		Tags:            request.Tags,
	}

	achievement, err := s.repoAchievement.Create(ctx, *achievement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data":   err.Error(),
		})
	}

	response := model.WebResponse[model.AchievementMongo]{
		Status: "success",
		Data:   *achievement,
	}

	return c.Status(fiber.StatusCreated).JSON(response)

}

func (s *AchievementServiceImpl) Update(c *fiber.Ctx) error {
	Id := c.Params("id")
	var request model.UpdateAchievementRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"errors": err.Error(),
		})
	}
	request.ID = Id
	ctx := c.UserContext()
	val := ctx.Value("user")
	achievement := &model.AchievementMongo{
		StudentID:       val.(*model.Claims).UserID,
		AchievementType: request.AchievementType,
		Title:           request.Title,
		Description:     request.Description,
		Details:         request.Details,
		Tags:            request.Tags,
	}

	achievement, err := s.repoAchievement.Update(ctx, *achievement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"errors": err.Error(),
		})
	}

	response := model.WebResponse[model.AchievementMongo]{
		Status: "success",
		Data:   *achievement,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *AchievementServiceImpl) Delete(c *fiber.Ctx) error {
	//TODO implement me
	Id := c.Params("id")
	ctx := c.UserContext()

	panic("implement me")
}

func (s *AchievementServiceImpl) FindByID(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *AchievementServiceImpl) FindAll(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *AchievementServiceImpl) Verify(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *AchievementServiceImpl) Submit(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *AchievementServiceImpl) History(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *AchievementServiceImpl) Attachment(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
