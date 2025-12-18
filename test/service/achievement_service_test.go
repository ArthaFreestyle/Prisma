package service_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"prisma/app/model"
	"prisma/app/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- MOCK REPOSITORIES ---

// 1. Mock Student Repository
type MockStudentRepo struct {
	mock.Mock
}

// Implementasi method yang dipanggil di function Create
func (m *MockStudentRepo) FindByUserId(ctx context.Context, userid string) (*model.Student, error) {
	args := m.Called(ctx, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

// Stub method lain untuk memenuhi interface StudentRepository
func (m *MockStudentRepo) Save(ctx context.Context, tx *sql.Tx, Student *model.Student) (*model.Student, error) {
	return nil, nil
}
func (m *MockStudentRepo) FindAll(ctx context.Context) ([]model.UserProfile, error) { return nil, nil }
func (m *MockStudentRepo) FindById(ctx context.Context, id string) (*model.UserProfile, error) {
	return nil, nil
}
func (m *MockStudentRepo) DeleteById(ctx context.Context, tx *sql.Tx, id string) error { return nil }
func (m *MockStudentRepo) UpdateById(ctx context.Context, Student *model.Student) (*model.Student, error) {
	return nil, nil
}

// 2. Mock Achievement Repository (Mongo)
type MockAchievementRepo struct {
	mock.Mock
}

func (m *MockAchievementRepo) Create(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error) {
	args := m.Called(ctx, Achievement)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AchievementMongo), args.Error(1)
}

// Stub method lain
func (m *MockAchievementRepo) Update(ctx context.Context, Achievement model.AchievementMongo) (*model.AchievementMongo, error) {
	return nil, nil
}
func (m *MockAchievementRepo) FindAll(ctx context.Context, Id []string) ([]model.AchievementMongo, error) {
	return nil, nil
}
func (m *MockAchievementRepo) FindById(ctx context.Context, id string) (*model.AchievementMongo, error) {
	return nil, nil
}

// 3. Mock Reference Repository (Postgres)
type MockReferenceRepo struct {
	mock.Mock
}

func (m *MockReferenceRepo) Create(ctx context.Context, achievement model.AchievementReference) (*model.AchievementReference, error) {
	args := m.Called(ctx, achievement)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}

// Stub method lain
func (m *MockReferenceRepo) Update(ctx context.Context, achievement model.AchievementReference) (*model.AchievementReference, error) {
	return nil, nil
}
func (m *MockReferenceRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *MockReferenceRepo) FindByID(ctx context.Context, id string) (*model.AchievementReferenceDetail, error) {
	return nil, nil
}
func (m *MockReferenceRepo) FindByLecturer(ctx context.Context, id string, page int, limit int) ([]model.AchievementReferenceLecturer, error) {
	return nil, nil
}
func (m *MockReferenceRepo) FindByStudent(ctx context.Context, id string, page int, limit int) ([]model.AchievementReferenceStudent, error) {
	return nil, nil
}
func (m *MockReferenceRepo) FindAll(ctx context.Context, page int, limit int) ([]model.AchievementReferenceAdmin, error) {
	return nil, nil
}
func (m *MockReferenceRepo) FindByStudentId(ctx context.Context, id string, page int, limit int) ([]model.AchievementReferenceAdmin, error) {
	return nil, nil
}

// --- UNIT TEST FUNCTION ---

func TestAchievementServiceImpl_Create(t *testing.T) {
	// Setup Mocks
	mockStudentRepo := new(MockStudentRepo)
	mockAchievementRepo := new(MockAchievementRepo)
	mockRefRepo := new(MockReferenceRepo)
	validator := validator.New()
	logger := logrus.New()

	// Inisialisasi Service
	svc := service.NewAchievementService(
		mockAchievementRepo,
		mockStudentRepo,
		mockRefRepo,
		validator,
		logger,
	)

	// Setup Fiber App
	app := fiber.New()

	// Middleware Mock Auth (Inject User Claims ke Context)
	app.Use(func(c *fiber.Ctx) error {
		// Mock Claims sesuai struct di code kamu
		claims := &model.Claims{
			UserID: "user-123",
			Role:   "mahasiswa",
		}
		// Inject ke UserContext karena di service pakai c.UserContext().Value("user")
		ctx := context.WithValue(c.UserContext(), "user", claims)
		c.SetUserContext(ctx)
		return c.Next()
	})

	// Register Route
	app.Post("/achievements", svc.Create)

	// Skenario Test
	t.Run("Success Create Achievement", func(t *testing.T) {
		// 1. Data Mock Return
		studentMock := &model.Student{
			ID:     "student-id-1",
			UserID: "user-123",
		}
		mongoID := primitive.NewObjectID()
		mongoResult := &model.AchievementMongo{
			ID:        mongoID,
			Title:     "Lomba Coding",
			StudentID: "student-id-1",
		}
		refResult := &model.AchievementReference{
			ID:                 "ref-id-1",
			StudentID:          "student-id-1",
			MongoAchievementID: mongoID.Hex(),
			Status:             "draft",
		}

		// 2. Expectation (Apa yang diharapkan dipanggil)
		mockStudentRepo.On("FindByUserId", mock.Anything, "user-123").Return(studentMock, nil)
		mockAchievementRepo.On("Create", mock.Anything, mock.MatchedBy(func(arg model.AchievementMongo) bool {
			return arg.Title == "Lomba Coding" && arg.StudentID == "student-id-1"
		})).Return(mongoResult, nil)
		mockRefRepo.On("Create", mock.Anything, mock.MatchedBy(func(arg model.AchievementReference) bool {
			return arg.MongoAchievementID == mongoID.Hex() && arg.Status == "draft"
		})).Return(refResult, nil)

		// 3. Request Payload
		payload := model.CreateAchievementRequest{
			AchievementType: "competition",
			Title:           "Lomba Coding",
			Description:     "Juara 1",
			Tags:            []string{"coding"},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/achievements", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// 4. Execute Request
		resp, err := app.Test(req)

		// 5. Assertions
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		// Verifikasi Mock dipanggil
		mockStudentRepo.AssertExpectations(t)
		mockAchievementRepo.AssertExpectations(t)
		mockRefRepo.AssertExpectations(t)
	})

	t.Run("Error Validation Failed", func(t *testing.T) {
		// Payload kosong (Title required)
		payload := model.CreateAchievementRequest{
			AchievementType: "competition",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/achievements", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Error Student Not Found", func(t *testing.T) {
		// Reset mock expectation
		mockStudentRepo.ExpectedCalls = nil

		// Expect FindByUserId return Error
		mockStudentRepo.On("FindByUserId", mock.Anything, "user-123").Return(nil, errors.New("record not found"))

		payload := model.CreateAchievementRequest{
			AchievementType: "competition",
			Title:           "Lomba Invalid",
			Description:     "Test",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/achievements", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Error Database Mongo Failed", func(t *testing.T) {
		mockStudentRepo.ExpectedCalls = nil
		mockAchievementRepo.ExpectedCalls = nil

		studentMock := &model.Student{ID: "student-id-1"}

		// Student ketemu, tapi Mongo Error
		mockStudentRepo.On("FindByUserId", mock.Anything, "user-123").Return(studentMock, nil)
		mockAchievementRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("mongo connection error"))

		payload := model.CreateAchievementRequest{
			AchievementType: "competition",
			Title:           "Lomba Mongo Error",
			Description:     "Test",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/achievements", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
