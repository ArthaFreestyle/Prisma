package service

import (
	"fmt"
	"os"
	"prisma/app/model"
	"prisma/app/repository"
	"sort"
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
	Reject(c *fiber.Ctx) error
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

// Create godoc
// @Summary      Create a new achievement
// @Description  Create a new achievement draft for a student.
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        request body model.CreateAchievementRequest true "Create Request"
// @Success      201  {object}  model.WebResponse[model.AchievementReference]
// @Failure      400  {object}  model.WebResponse[string]
// @Failure      404  {object}  model.WebResponse[string]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements [post]
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

// Update godoc
// @Summary      Update an achievement
// @Description  Update achievement details. Students can only update if status is 'draft'.
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Achievement ID"
// @Param        request body model.UpdateAchievementRequest true "Update Request"
// @Success      200  {object}  model.WebResponse[model.AchievementReferenceDetail]
// @Failure      400  {object}  model.WebResponse[string]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements/{id} [put]
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

	if val.(*model.Claims).Role == "mahasiswa" && Achievement.Status != "draft" {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: fmt.Sprintf("Achievement %s is already submitted", request.ID),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
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

// Delete godoc
// @Summary      Delete an achievement
// @Description  Delete achievement by ID.
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Achievement ID"
// @Success      200  {object}  model.WebResponse[string]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements/{id} [delete]
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

// FindByID godoc
// @Summary      Get achievement by ID
// @Description  Retrieve full details of an achievement including mongo data.
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Achievement ID"
// @Success      200  {object}  model.WebResponse[model.AchievementReferenceDetail]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements/{id} [get]
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

// FindAll godoc
// @Summary      Get all achievements
// @Description  Get achievements with pagination. Results filtered by user role (admin, mahasiswa, lecturer).
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Limit per page"
// @Success      200  {object}  model.WebResponse[[]model.AchievementReferenceDetail]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements [get]
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

// Verify godoc
// @Summary      Verify an achievement
// @Description  Mark an achievement as verified.
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Achievement ID"
// @Success      200  {object}  model.WebResponse[model.AchievementReferenceDetail]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements/{id}/verify [patch]
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
	Achievement.Status = "verified"
	Achievement.VerifiedBy = &val.(*model.Claims).UserID
	Achievement.VerifiedAt = &now

	response := model.WebResponse[*model.AchievementReferenceDetail]{
		Status: "success",
		Data:   Achievement,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Submit godoc
// @Summary      Submit an achievement
// @Description  Submit an achievement for review.
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Achievement ID"
// @Success      200  {object}  model.WebResponse[model.AchievementReferenceDetail]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements/{id}/submit [patch]
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

// History godoc
// @Summary      Get achievement history
// @Description  Get the timeline history of an achievement.
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Achievement ID"
// @Success      200  {object}  model.WebResponse[[]model.AchievementHistory]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements/{id}/history [get]
func (s *AchievementServiceImpl) History(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()

	achievement, err := s.repoAchivementReference.FindByID(ctx, id)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	histories := make([]model.AchievementHistory, 0, 3)

	histories = append(histories, model.AchievementHistory{
		Action:    "Dibuat",
		Timestamp: achievement.CreatedAt,
	})

	if achievement.SubmittedAt != nil {
		histories = append(histories, model.AchievementHistory{
			Action:    "Diajukan",
			Timestamp: *achievement.SubmittedAt,
		})
	}

	if achievement.VerifiedAt != nil {
		actionLabel := "Diverifikasi"
		if achievement.Status == "rejected" {
			actionLabel = "rejected"
		} else if achievement.Status == "verified" {
			actionLabel = "verified"
		}

		histories = append(histories, model.AchievementHistory{
			Action:    actionLabel,
			Timestamp: *achievement.VerifiedAt,
		})
	}

	sort.Slice(histories, func(i, j int) bool {
		return histories[i].Timestamp.After(histories[j].Timestamp)
	})

	return c.JSON(model.WebResponse[[]model.AchievementHistory]{
		Status: "success",
		Data:   histories,
	})
}

// Attachment godoc
// @Summary      Upload attachment
// @Description  Upload file attachments for an achievement.
// @Tags         Achievement
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   path      string  true  "Achievement ID"
// @Param        attachments formData file true "Files to upload"
// @Success      200  {object}  model.WebResponse[model.AchievementReferenceDetail]
// @Failure      400  {object}  model.WebResponse[string]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements/{id}/attachments [post]
func (s *AchievementServiceImpl) Attachment(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.WebResponse[string]{
			Status: "error",
			Errors: "Gagal memproses form upload",
		})
	}

	files := form.File["attachments"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.WebResponse[string]{
			Status: "error",
			Errors: "Tidak ada file yang diupload",
		})
	}

	achievementRef, err := s.repoAchivementReference.FindByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[string]{
			Status: "error",
			Errors: "Ref not found: " + err.Error(),
		})
	}

	achievementObj, err := s.repoAchievement.FindById(ctx, achievementRef.MongoAchievementID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[string]{
			Status: "error",
			Errors: "Detail not found: " + err.Error(),
		})
	}

	var newAttachments []model.Attachment

	baseDir := "./public/uploads/achievements"
	baseURL := "/uploads/achievements"

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		os.MkdirAll(baseDir, 0755)
	}

	for _, file := range files {
		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		savePath := fmt.Sprintf("%s/%s", baseDir, uniqueName)

		if err := c.SaveFile(file, savePath); err != nil {
			response := model.WebResponse[string]{
				Status: "error",
				Errors: err.Error(),
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		attachment := model.Attachment{
			FileName:   file.Filename,
			FileURL:    fmt.Sprintf("%s/%s", baseURL, uniqueName),
			FileType:   file.Header.Get("Content-Type"),
			UploadedAt: time.Now(),
		}

		newAttachments = append(newAttachments, attachment)
	}

	achievementObj.Attachments = append(achievementObj.Attachments, newAttachments...)

	AchievementObj, err := s.repoAchievement.Update(ctx, *achievementObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[string]{
			Status: "error",
			Errors: "Gagal update database: " + err.Error(),
		})
	}

	achievementRef.Detail = AchievementObj

	return c.JSON(model.WebResponse[*model.AchievementReferenceDetail]{
		Status: "success",
		Data:   achievementRef,
	})
}

// Reject godoc
// @Summary      Reject an achievement
// @Description  Reject a submitted achievement with a note.
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Achievement ID"
// @Param        request body model.CreateRejection true "Rejection Note"
// @Success      200  {object}  model.WebResponse[model.AchievementReferenceDetail]
// @Failure      400  {object}  model.WebResponse[string]
// @Failure      500  {object}  model.WebResponse[string]
// @Security     BearerAuth
// @Router       /achievements/{id}/reject [post]
func (s *AchievementServiceImpl) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()
	val := ctx.Value("user")

	var request model.CreateRejection
	if err := c.BodyParser(&request); err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Errors: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

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
		Status:             "rejected",
		RejectionNote:      request.RejectionNote,
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
