package service

import (
	"prisma/app/model"
	"prisma/app/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	repoAchievement         repository.AchievementRepository
	repoStudent             repository.StudentRepository
	repoAchivementReference repository.AchievementReferenceRepository
	validate                *validator.Validate
	Log                     *logrus.Logger
}

func NewAchievementService(repo repository.AchievementRepository, repoStudent repository.StudentRepository, repoAchievementReference repository.AchievementReferenceRepository, validate *validator.Validate, Log *logrus.Logger) *AchievementServiceImpl {
	return &AchievementServiceImpl{
		repoAchievement:         repo,
		validate:                validate,
		repoStudent:             repoStudent,
		repoAchivementReference: repoAchievementReference,
		Log:                     Log,
	}
}

func (s *AchievementServiceImpl) Create(c *fiber.Ctx) error {
	var request model.CreateAchievementRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": err.Error()})
	}

	if err := s.validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": err.Error()})
	}

	ctx := c.UserContext()
	val := ctx.Value("user")

	Student, err := s.repoStudent.FindByUserId(ctx, val.(*model.Claims).UserID)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	achievement := &model.AchievementMongo{
		StudentID:       Student.ID, // Pakai ID Student, bukan User ID
		AchievementType: request.AchievementType,
		Title:           request.Title,
		Description:     request.Description,
		Details:         request.Details,
		Tags:            request.Tags,
	}

	createdMongo, err := s.repoAchievement.Create(ctx, *achievement)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	ref := model.AchievementReference{
		StudentID:          Student.ID,
		MongoAchievementID: createdMongo.ID.Hex(),
		Status:             "draft",
	}

	createdRef, err := s.repoAchivementReference.Create(ctx, ref)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	createdRef.Detail = createdMongo

	return c.Status(fiber.StatusCreated).JSON(model.WebResponse[model.AchievementReference]{
		Status: "success",
		Data:   *createdRef,
	})
}

func (s *AchievementServiceImpl) Update(c *fiber.Ctx) error {
	Id := c.Params("id")
	var request model.UpdateAchievementRequest
	if err := c.BodyParser(&request); err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	request.ID = Id
	ctx := c.UserContext()
	val := ctx.Value("user")
	Achievement, err := s.repoAchivementReference.FindByID(ctx, request.ID)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	id, err := primitive.ObjectIDFromHex(Achievement.ID)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	achievementObj := &model.AchievementMongo{
		ID:              id,
		StudentID:       val.(*model.Claims).UserID,
		AchievementType: request.AchievementType,
		Title:           request.Title,
		Description:     request.Description,
		Details:         request.Details,
		Tags:            request.Tags,
	}

	achievementObj, err = s.repoAchievement.Update(ctx, *achievementObj)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	Achievement.Detail = achievementObj

	response := model.WebResponse[model.AchievementReferenceDetail]{
		Status: "success",
		Data:   *Achievement,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *AchievementServiceImpl) Delete(c *fiber.Ctx) error {
	//TODO implement me
	Id := c.Params("id")
	ctx := c.UserContext()

	err := s.repoAchivementReference.Delete(ctx, Id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	response := model.WebResponse[string]{
		Status: "success",
		Data:   "Data has been deleted",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *AchievementServiceImpl) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()

	Achievement, err := s.repoAchivementReference.FindByID(ctx, id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	AchievementObj, err := s.repoAchievement.FindById(ctx, Achievement.MongoAchievementID)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	Achievement.Detail = AchievementObj

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[*model.AchievementReferenceDetail]{
		Status: "success",
		Data:   Achievement,
	})
}

func (s *AchievementServiceImpl) FindAll(c *fiber.Ctx) error {

	Page := c.QueryInt("page", 1)
	Limit := c.QueryInt("limit", 10)
	ctx := c.UserContext()
	val := ctx.Value("user")
	var response model.WebResponse[any]
	if val.(*model.Claims).Role == "admin" {
		Achievements, err := s.repoAchivementReference.FindAll(ctx, Page, Limit)
		if err != nil {
			response := model.WebResponse[string]{
				Status: "error",
				Errors: err.Error(),
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
		var oids []string
		for _, ach := range Achievements {
			oids = append(oids, ach.MongoAchievementID)
		}

		if len(oids) == 0 {
			response.Data = []any{} // Return array kosong
			return c.Status(fiber.StatusOK).JSON(response)
		}
		AchievementObjs, err := s.repoAchievement.FindAll(ctx, oids)
		if err != nil {
			response := model.WebResponse[string]{
				Status: "error",
				Errors: err.Error(),
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		detailMap := make(map[primitive.ObjectID]model.AchievementMongo)
		for _, detail := range AchievementObjs {
			detailMap[detail.ID] = detail
		}

		for i := range Achievements {
			mongoID, err := primitive.ObjectIDFromHex(Achievements[i].MongoAchievementID)

			if err == nil {
				if detail, found := detailMap[mongoID]; found {
					Achievements[i].Title = detail.Title
					Achievements[i].Type = detail.AchievementType
					Achievements[i].CreatedAt = detail.CreatedAt

				}
			}
		}
		response.Data = Achievements
	} else if val.(*model.Claims).Role == "mahasiswa" {
		Achievements, err := s.repoAchivementReference.FindByStudent(ctx, val.(*model.Claims).UserID, Page, Limit)
		if err != nil {
			response := model.WebResponse[string]{
				Status: "error",
				Errors: err.Error(),
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
		var oids []string
		for _, ach := range Achievements {
			oids = append(oids, ach.MongoAchievementID)
		}

		if len(oids) == 0 {
			response.Data = []any{} // Return array kosong
			return c.Status(fiber.StatusOK).JSON(response)
		}
		AchievementObjs, err := s.repoAchievement.FindAll(ctx, oids)
		if err != nil {
			response := model.WebResponse[string]{
				Status: "error",
				Errors: err.Error(),
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		detailMap := make(map[primitive.ObjectID]model.AchievementMongo)
		for _, detail := range AchievementObjs {
			detailMap[detail.ID] = detail
		}

		for i := range Achievements {
			mongoID, err := primitive.ObjectIDFromHex(Achievements[i].MongoAchievementID)

			if err == nil {
				if detail, found := detailMap[mongoID]; found {
					Achievements[i].Title = detail.Title
					Achievements[i].Type = detail.AchievementType
					Achievements[i].CreatedAt = detail.CreatedAt

				}
			}
		}
		response.Data = Achievements
	} else if val.(*model.Claims).Role == "lecturer" {
		Achievements, err := s.repoAchivementReference.FindByLecturer(ctx, val.(*model.Claims).UserID, Page, Limit)
		if err != nil {
			response := model.WebResponse[string]{
				Status: "error",
				Errors: err.Error(),
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
		var oids []string
		for _, ach := range Achievements {
			oids = append(oids, ach.MongoAchievementID)
		}

		if len(oids) == 0 {
			response.Data = []any{} // Return array kosong
			return c.Status(fiber.StatusOK).JSON(response)
		}
		AchievementObjs, err := s.repoAchievement.FindAll(ctx, oids)
		if err != nil {
			response := model.WebResponse[string]{
				Status: "error",
				Errors: err.Error(),
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		detailMap := make(map[primitive.ObjectID]model.AchievementMongo)
		for _, detail := range AchievementObjs {
			detailMap[detail.ID] = detail
		}

		for i := range Achievements {
			mongoID, err := primitive.ObjectIDFromHex(Achievements[i].MongoAchievementID)

			if err == nil {
				if detail, found := detailMap[mongoID]; found {
					Achievements[i].Title = detail.Title
					Achievements[i].Type = detail.AchievementType
					Achievements[i].CreatedAt = detail.CreatedAt

				}
			}
		}
		response.Data = Achievements
	}
	response.Status = "success"
	response.Paging = &model.PageMetaData{
		Page: Page,
		Size: Limit,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *AchievementServiceImpl) Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()
	val := ctx.Value("user")
	Achievement, err := s.repoAchivementReference.FindByID(ctx, id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	now := time.Now()
	AchievementRefer := &model.AchievementReference{
		MongoAchievementID: Achievement.MongoAchievementID,
		StudentID:          Achievement.UserDetail.StudentProfile.StudentID,
		ID:                 Achievement.ID,
		Status:             "verified",
		SubmittedAt:        Achievement.SubmittedAt,
		VerifiedAt:         &now,
		VerifiedBy:         val.(*model.Claims).UserID,
		Detail:             Achievement.Detail,
	}

	AchievementRefer, err = s.repoAchivementReference.Update(ctx, *AchievementRefer)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	Achievement.Status = AchievementRefer.Status
	Achievement.VerifiedAt = Achievement.VerifiedAt
	Achievement.VerifiedBy = Achievement.VerifiedBy

	response := model.WebResponse[*model.AchievementReferenceDetail]{
		Status: "success",
		Data:   Achievement,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *AchievementServiceImpl) Submit(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()
	val := ctx.Value("user")

	Achievement, err := s.repoAchivementReference.FindByID(ctx, id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	now := time.Now()
	AchievementRefer := &model.AchievementReference{
		MongoAchievementID: Achievement.MongoAchievementID,
		StudentID:          Achievement.UserDetail.StudentProfile.StudentID,
		ID:                 Achievement.ID,
		Status:             "submitted",
		SubmittedAt:        Achievement.SubmittedAt,
		VerifiedAt:         &now,
		VerifiedBy:         val.(*model.Claims).UserID,
		Detail:             Achievement.Detail,
	}

	AchievementRefer, err = s.repoAchivementReference.Update(ctx, *AchievementRefer)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	Achievement.Status = AchievementRefer.Status

	response := model.WebResponse[*model.AchievementReferenceDetail]{
		Status: "success",
		Data:   Achievement,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *AchievementServiceImpl) History(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *AchievementServiceImpl) Attachment(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
