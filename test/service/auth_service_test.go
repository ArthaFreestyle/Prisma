package service_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"prisma/app/model"
	"prisma/app/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// --- MOCK REPOSITORIES ---

// 1. Mock User Repository
type MockUserRepoAuth struct {
	mock.Mock
}

func (m *MockUserRepoAuth) FindByUsername(ctx context.Context, Username string) (*model.User, error) {
	args := m.Called(ctx, Username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// Stub method lain
func (m *MockUserRepoAuth) Save(ctx context.Context, tx *sql.Tx, User *model.User) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepoAuth) Update(ctx context.Context, User model.User) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepoAuth) UpdateRole(ctx context.Context, tx *sql.Tx, User model.User) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepoAuth) Delete(ctx context.Context, UserId string) error { return nil }
func (m *MockUserRepoAuth) FindById(ctx context.Context, UserId string) (*model.UserProfile, error) {
	return nil, nil
}
func (m *MockUserRepoAuth) FindAll(ctx context.Context) (*[]model.User, error) { return nil, nil }

// 2. Mock Auth Repository (Redis)
type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) Logout(ctx context.Context, RefreshToken string) error {
	args := m.Called(ctx, RefreshToken)
	return args.Error(0)
}

func (m *MockAuthRepo) RefreshToken(ctx context.Context, RefreshToken string, jwtKey []byte) (string, error) {
	args := m.Called(ctx, RefreshToken, jwtKey)
	return args.String(0), args.Error(1)
}

// --- HELPER FOR TOKEN ---
func generateValidRefreshToken(secret []byte) string {
	// Membuat token dummy yang valid secara struktur JWT agar lolos utils.ValidateToken
	claims := jwt.MapClaims{
		"exp":      time.Now().Add(time.Hour).Unix(),
		"user_id":  "user-123",
		"username": "testuser",
		"role":     "mahasiswa",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString(secret)
	return t
}

// --- UNIT TEST FUNCTION ---

func TestAuthServiceImpl(t *testing.T) {
	// Setup Dependencies
	mockUserRepo := new(MockUserRepoAuth)
	mockAuthRepo := new(MockAuthRepo)
	logger := logrus.New()
	secretKey := []byte("secret-key-test") // Secret key dummy

	// Init Service
	svc := service.NewAuthService(
		mockUserRepo,
		mockAuthRepo,
		logger,
		secretKey,
	)

	app := fiber.New()

	// Setup Routes
	app.Post("/auth/login", svc.Login)
	app.Post("/auth/logout", svc.Logout)
	app.Post("/auth/refresh", svc.RefreshToken)

	t.Run("Login Success", func(t *testing.T) {
		// Arrange
		password := "password123"
		// Kita harus generate hash asli supaya utils.CheckPasswordHash berhasil
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		userMock := &model.User{
			ID:           "user-1",
			Username:     "testuser",
			PasswordHash: string(hashedPassword),
			RoleName:     "mahasiswa",
		}

		mockUserRepo.On("FindByUsername", mock.Anything, "testuser").Return(userMock, nil)

		payload := model.LoginRequest{
			Username: "testuser",
			Password: password,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Cek apakah response mengandung token
		var respBody model.WebResponse[model.LoginResponse]
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NotEmpty(t, respBody.Data.Token)
		assert.NotEmpty(t, respBody.Data.RefreshToken)
	})

	t.Run("Login Wrong Password", func(t *testing.T) {
		// Arrange
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
		userMock := &model.User{
			Username:     "testuser",
			PasswordHash: string(hashedPassword),
		}

		// User ditemukan, tapi password yang dikirim nanti salah
		mockUserRepo.On("FindByUsername", mock.Anything, "testuser").Return(userMock, nil)

		payload := model.LoginRequest{
			Username: "testuser",
			Password: "wrong-password",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Act
		resp, err := app.Test(req)

		// Assert
		// Fiber.ErrUnauthorized biasanya map ke 401
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Login User Not Found", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserRepo.On("FindByUsername", mock.Anything, "unknown").Return(nil, errors.New("user not found"))

		payload := model.LoginRequest{Username: "unknown", Password: "pwd"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Logout Success", func(t *testing.T) {
		// Arrange
		mockAuthRepo.On("Logout", mock.Anything, "dummy-refresh-token").Return(nil)

		req := httptest.NewRequest("POST", "/auth/logout", nil)
		// Set Cookie
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "dummy-refresh-token"})

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockAuthRepo.AssertExpectations(t)
	})

	t.Run("RefreshToken Success", func(t *testing.T) {
		// Arrange
		// Kita butuh token valid karena service memanggil utils.ValidateToken
		validToken := generateValidRefreshToken(secretKey)
		newAccessToken := "new-access-token-from-redis"

		// Mock Auth Repo harus dipanggil
		mockAuthRepo.On("RefreshToken", mock.Anything, validToken, secretKey).Return(newAccessToken, nil)

		req := httptest.NewRequest("POST", "/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: validToken})

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody model.WebResponse[model.LoginResponse]
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, newAccessToken, respBody.Data.Token)
	})

	t.Run("RefreshToken Invalid Token", func(t *testing.T) {
		// Arrange: Token sembarangan yang tidak valid signature-nya
		invalidToken := "invalid.jwt.token"

		// Repo tidak akan dipanggil karena validasi token gagal duluan
		mockAuthRepo.ExpectedCalls = nil

		req := httptest.NewRequest("POST", "/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: invalidToken})

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}
