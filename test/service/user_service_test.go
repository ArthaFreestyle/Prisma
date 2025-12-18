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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCKS FOR REPOSITORIES ---

// 1. Mock User Repository
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Save(ctx context.Context, tx *sql.Tx, User *model.User) (*model.User, error) {
	// Kita gunakan mock.Anything untuk tx karena tx dibuat di dalam service (sulit di-match exact object-nya)
	args := m.Called(ctx, tx, User)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// Stub method lain (biar implement interface)
func (m *MockUserRepo) Update(ctx context.Context, User model.User) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepo) UpdateRole(ctx context.Context, tx *sql.Tx, User model.User) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Delete(ctx context.Context, UserId string) error { return nil }
func (m *MockUserRepo) FindById(ctx context.Context, UserId string) (*model.UserProfile, error) {
	return nil, nil
}
func (m *MockUserRepo) FindAll(ctx context.Context) (*[]model.User, error) { return nil, nil }
func (m *MockUserRepo) FindByUsername(ctx context.Context, Username string) (*model.User, error) {
	return nil, nil
}

// 3. Mock Lecturer Repository
type MockLecturerRepo struct {
	mock.Mock
}

func (m *MockLecturerRepo) FindAll(ctx context.Context) ([]model.UserProfile, error) {
	return nil, nil
}

func (m *MockLecturerRepo) FindById(ctx context.Context, id string) (*model.UserProfile, error) {
	return nil, nil
}

func (m *MockLecturerRepo) FindAllAdvices(ctx context.Context, id string) ([]model.UserProfile, error) {
	return nil, nil
}

func (m *MockLecturerRepo) Save(ctx context.Context, tx *sql.Tx, Lecturer *model.Lecturer) (*model.Lecturer, error) {
	args := m.Called(ctx, tx, Lecturer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

// Stub method lain
func (m *MockLecturerRepo) DeleteById(ctx context.Context, tx *sql.Tx, id string) error { return nil }

// --- UNIT TEST FUNCTION ---

func TestUserServiceImpl_Create(t *testing.T) {
	// 1. Setup SQL Mock (Simulasi Database Connection & Transaction)
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// 2. Setup Dependencies
	mockUserRepo := new(MockUserRepo)
	mockStudentRepo := new(MockStudentRepo)
	mockLecturerRepo := new(MockLecturerRepo)
	validate := validator.New()
	logger := logrus.New()

	// Inisialisasi Service dengan SQL DB Mock
	svc := service.NewUserService(
		mockUserRepo,
		mockStudentRepo,
		mockLecturerRepo,
		db, // Inject DB mock disini
		validate,
		logger,
	)

	app := fiber.New()
	app.Post("/users", svc.Create)

	// Hardcoded Role IDs sesuai kode aslimu
	roleStudent := "11111111-1111-1111-1111-111111111111"
	// roleLecturer := "22222222-2222-2222-2222-222222222222"

	t.Run("Success Create User Student", func(t *testing.T) {
		// Arrange
		payload := model.UserCreateRequest{
			Username: "maba2024",
			Email:    "maba@univ.ac.id",
			Password: "password123", // Akan di hash di service
			FullName: "Mahasiswa Baru",
			RoleID:   roleStudent,
			StudentProfile: &model.StudentCreate{
				StudentID:    "NIM123",
				ProgramStudy: "Informatika",
				AcademicYear: "2024",
				AdvisorID:    "DOSEN001",
			},
		}

		// Mock Return Data
		createdUser := &model.User{
			ID:       "user-uuid",
			Username: payload.Username,
			RoleId:   payload.RoleID,
		}
		createdStudent := &model.Student{
			ID:        "student-uuid",
			UserID:    "user-uuid",
			StudentID: "NIM123",
		}

		// --- Expectations Sequence ---

		// 1. Service akan memulai Transaction
		sqlMock.ExpectBegin()

		// 2. User Repo Save dipanggil
		// Note: matchers untuk argument struct user agak longgar (mock.Anything) karena password sudah di-hash
		mockUserRepo.On("Save", mock.Anything, mock.AnythingOfType("*sql.Tx"), mock.MatchedBy(func(u *model.User) bool {
			return u.Username == "maba2024" && u.RoleId == roleStudent
		})).Return(createdUser, nil)

		// 3. Student Repo Save dipanggil (karena role = student)
		mockStudentRepo.On("Save", mock.Anything, mock.AnythingOfType("*sql.Tx"), mock.MatchedBy(func(s *model.Student) bool {
			return s.StudentID == "NIM123" && s.UserID == "user-uuid"
		})).Return(createdStudent, nil)

		// 4. Jika sukses semua, Transaction Commit
		sqlMock.ExpectCommit()

		// Act
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, sqlMock.ExpectationsWereMet()) // Verifikasi urutan SQL
		mockUserRepo.AssertExpectations(t)
		mockStudentRepo.AssertExpectations(t)
	})

	t.Run("Error Transaction Rollback on User Save Fail", func(t *testing.T) {
		// Arrange
		payload := model.UserCreateRequest{
			Username:       "failuser",
			Email:          "fail@test.com",
			Password:       "pass",
			FullName:       "Fail User",
			RoleID:         roleStudent,
			StudentProfile: &model.StudentCreate{StudentID: "1"},
		}

		// Expectations
		sqlMock.ExpectBegin() // 1. Begin

		// 2. Repo User Save Error (misal username duplicate)
		mockUserRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("username already exists"))

		sqlMock.ExpectRollback() // 3. HARUS Rollback

		// Act
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Error Validation Missing Field", func(t *testing.T) {
		// Arrange: Payload tidak lengkap (RoleID required misal)
		payload := model.UserCreateRequest{
			Username: "incomplete",
		}

		// Tidak ada interaksi database sama sekali karena validasi gagal duluan

		// Act
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
