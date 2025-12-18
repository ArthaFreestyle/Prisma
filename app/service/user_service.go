package service

import (
	"database/sql"
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
	UpdateRole(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	FindById(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	Profile(c *fiber.Ctx) error
}

func NewUserService(repoUser repository.UserRepository, repoStudent repository.StudentRepository, repoLecturer repository.LecturerRepository, DB *sql.DB, validate *validator.Validate, log *logrus.Logger) UserService {
	return &UserServiceImpl{repoUser, repoStudent, repoLecturer, DB, validate, log}
}

type UserServiceImpl struct {
	repoUser     repository.UserRepository
	repoStudent  repository.StudentRepository
	repoLecturer repository.LecturerRepository
	DB           *sql.DB
	validate     *validator.Validate
	Log          *logrus.Logger
}

// UpdateRole godoc
// @Summary Update user role
// @Description Update user role and associated profile data (student/lecturer)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body model.UserUpdateRole true "User role update request"
// @Success 200 {object} model.SwaggerWebResponseInterface "Successfully updated user role"
// @Failure 400 {object} model.SwaggerWebResponseString "Bad request - invalid input or role ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /users/{id}/role [put]
func (s *UserServiceImpl) UpdateRole(c *fiber.Ctx) error {
	UserId := c.Params("id")
	var request model.UserUpdateRole
	ctx := c.UserContext()
	if err := c.BodyParser(&request); err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Data:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	if err := validator.New().Struct(request); err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Data:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !utils.CheckRoleAccepted(request.RoleID) {
		response := model.WebResponse[string]{
			Status: "success",
			Data:   "id role tidak ditemukan",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userProfile, err := s.repoUser.FindById(ctx, UserId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data":   err.Error(),
		})
	}

	users := &model.User{
		ID:       userProfile.User.ID,
		RoleId:   request.RoleID,
		FullName: userProfile.User.FullName,
		Username: userProfile.User.Username,
		Email:    userProfile.User.Email,
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data":   err.Error(),
		})
	}
	defer tx.Rollback()
	users, err = s.repoUser.UpdateRole(ctx, tx, *users)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data":   err.Error(),
		})
	}

	if userProfile.StudentID.Valid {
		err := s.repoStudent.DeleteById(ctx, tx, userProfile.StudentID.String)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data":   err.Error(),
			})
		}
	} else if userProfile.LecturerID.Valid {
		err := s.repoLecturer.DeleteById(ctx, tx, userProfile.LecturerID.String)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data":   err.Error(),
			})
		}
	}

	var UserData interface{}
	switch request.RoleID {
	case "11111111-1111-1111-1111-111111111111":
		student := &model.Student{
			UserID:       users.ID,
			StudentID:    request.StudentData.StudentID,
			ProgramStudy: request.StudentData.ProgramStudy,
			AcademicYear: request.StudentData.AcademicYear,
			AdvisorID:    request.StudentData.AdvisorID,
		}
		student, err = s.repoStudent.Save(ctx, tx, student)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data":   err.Error(),
			})
		}
		UserData = model.UserCreateResponse{
			ID:             users.ID,
			Username:       users.Username,
			Email:          users.Email,
			FullName:       users.FullName,
			RoleID:         users.RoleId,
			StudentProfile: student,
		}

	case "22222222-2222-2222-2222-222222222222":
		lecturer := &model.Lecturer{
			UserID:     users.ID,
			LecturerID: request.LecturerData.LecturerID,
			Department: request.LecturerData.Department,
		}
		lecturer, err := s.repoLecturer.Save(ctx, tx, lecturer)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data":   err.Error(),
			})
		}
		UserData = model.UserCreateResponse{
			ID:              users.ID,
			Username:        users.Username,
			Email:           users.Email,
			FullName:        users.FullName,
			RoleID:          users.RoleId,
			LecturerProfile: lecturer,
		}

	case "33333333-3333-3333-3333-333333333333":
		UserData = model.UserCreateResponse{
			ID:       users.ID,
			Username: users.Username,
			Email:    users.Email,
			FullName: users.FullName,
			RoleID:   users.RoleId,
		}

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "Role ID mismatch in processing"})
	}

	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data":   "Gagal menyimpan perubahan permanen: " + err.Error(),
		})
	}

	response := model.WebResponse[interface{}]{
		Status: "success",
		Data:   UserData,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// Create godoc
// @Summary Create new user
// @Description Create a new user with role-specific profile (student/lecturer/admin)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body model.UserCreateRequest true "User creation request"
// @Success 200 {object} model.SwaggerWebResponseInterface "Successfully created user"
// @Failure 400 {object} map[string]interface{} "Bad request - validation error or missing profile data"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /users [post]
func (s *UserServiceImpl) Create(c *fiber.Ctx) error {
	var request model.UserCreateRequest
	if err := c.BodyParser(&request); err != nil {
		s.Log.Info(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := s.validate.Struct(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	ctx := c.UserContext()
	tx, err := s.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data":   err.Error(),
		})
	}

	defer tx.Rollback()
	PasswordHash, err := utils.HashPassword(request.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user := &model.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: PasswordHash,
		FullName:     request.FullName,
		RoleId:       request.RoleID,
	}
	user, err = s.repoUser.Save(ctx, tx, user)
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !utils.CheckRoleAccepted(request.RoleID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"data":   "role not accepted",
		})
	}

	var UserData interface{}
	if request.RoleID == "11111111-1111-1111-1111-111111111111" {
		if request.StudentProfile == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"data":   "data student profile is nil",
			})
		}
		student := &model.Student{
			UserID:       user.ID,
			StudentID:    request.StudentProfile.StudentID,
			ProgramStudy: request.StudentProfile.ProgramStudy,
			AcademicYear: request.StudentProfile.AcademicYear,
			AdvisorID:    request.StudentProfile.AdvisorID,
		}
		student, err = s.repoStudent.Save(ctx, tx, student)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		UserData = model.UserCreateResponse{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			FullName:       user.FullName,
			RoleID:         user.RoleId,
			StudentProfile: student,
		}

	} else if request.RoleID == "22222222-2222-2222-2222-222222222222" {
		if request.LecturerProfile == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"data":   "tidak ada data lecturer",
			})
		}
		lecturer := &model.Lecturer{
			UserID:     user.ID,
			LecturerID: request.LecturerProfile.LecturerID,
			Department: request.LecturerProfile.Department,
		}
		lecturer, err := s.repoLecturer.Save(ctx, tx, lecturer)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		UserData = model.UserCreateResponse{
			ID:              user.ID,
			Username:        user.Username,
			Email:           user.Email,
			FullName:        user.FullName,
			RoleID:          user.RoleId,
			LecturerProfile: lecturer,
		}

	} else if request.RoleID == "33333333-3333-3333-3333-333333333333" {
		UserData = model.UserCreateResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			RoleID:   user.RoleId,
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"data":   "Gagal menyimpan perubahan permanen: " + err.Error(),
		})
	}
	response := model.WebResponse[interface{}]{
		Status: "success",
		Data:   UserData,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// Update godoc
// @Summary Update user information
// @Description Update user's basic information (username, email, fullname)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body model.UserCreateRequest true "User update request"
// @Success 200 {object} model.SwaggerWebResponseUserUpdateResponse "Successfully updated user"
// @Failure 400 {object} model.SwaggerWebResponseString "Bad request - invalid input"
// @Failure 404 {object} model.SwaggerWebResponseString "User not found"
// @Security BearerAuth
// @Router /users/{id} [put]
func (s *UserServiceImpl) Update(c *fiber.Ctx) error {
	UserId := c.Params("id")
	var request model.UserCreateRequest
	ctx := c.UserContext()
	if err := c.BodyParser(&request); err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Data:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	user := &model.User{
		ID:       UserId,
		Username: request.Username,
		Email:    request.Email,
		FullName: request.FullName,
	}

	user, err := s.repoUser.Update(ctx, *user)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Data:   err.Error(),
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	UserData := model.UserUpdateResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	}
	response := model.WebResponse[model.UserUpdateResponse]{
		Status: "success",
		Data:   UserData,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Delete godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.SwaggerWebResponseString "Successfully deleted user"
// @Failure 400 {object} map[string]interface{} "Bad request - user not found"
// @Security BearerAuth
// @Router /users/{id} [delete]
func (s *UserServiceImpl) Delete(c *fiber.Ctx) error {
	UserId := c.Params("id")
	ctx := c.UserContext()

	err := s.repoUser.Delete(ctx, UserId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	response := model.WebResponse[string]{
		Status: "success",
		Data:   "user deleted",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// FindById godoc
// @Summary Get user by ID
// @Description Get user details by ID including role-specific profile
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.SwaggerWebResponseUserResponse "Successfully retrieved user"
// @Failure 404 {object} model.SwaggerWebResponseString "User not found"
// @Security BearerAuth
// @Router /users/{id} [get]
func (s *UserServiceImpl) FindById(c *fiber.Ctx) error {
	var UserID = c.Params("id")
	ctx := c.UserContext()
	Users, err := s.repoUser.FindById(ctx, UserID)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Data:   "User not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}
	UserResponse := model.UserResponse{
		ID:       Users.User.ID,
		Username: Users.User.Username,
		Email:    Users.User.Email,
		FullName: Users.User.FullName,
		Role:     Users.User.RoleName,
	}

	if Users.StudentID.Valid {
		UserResponse.StudentProfile = &model.StudentCreate{
			StudentID:    Users.StudentID.String,
			ProgramStudy: Users.ProgramStudy.String,
			AcademicYear: Users.AcademicYear.String,
			AdvisorID:    Users.AdvisorID.String,
		}
	} else if Users.LecturerID.Valid {
		UserResponse.LecturerProfile = &model.LecturerCreate{
			LecturerID: Users.LecturerID.String,
			Department: Users.Department.String,
		}
	}

	response := model.WebResponse[model.UserResponse]{
		Status: "success",
		Data:   UserResponse,
	}

	return c.Status(fiber.StatusOK).JSON(response)

}

// FindAll godoc
// @Summary Get all users
// @Description Get list of all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} model.SwaggerWebResponseUserResponses "Successfully retrieved users"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /users [get]
func (s *UserServiceImpl) FindAll(c *fiber.Ctx) error {

	ctx := c.UserContext()
	Users, err := s.repoUser.FindAll(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var userResponses []model.UserResponse
	for _, u := range *Users {
		userResponses = append(userResponses, model.UserResponse{
			ID:       u.ID,
			Email:    u.Email,
			Username: u.Username,
			FullName: u.FullName,
			Role:     u.RoleName,
		})
	}
	response := model.WebResponse[[]model.UserResponse]{
		Status: "success",
		Data:   userResponses,
	}
	return c.Status(fiber.StatusOK).JSON(response)

}

// Profile godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} model.SwaggerWebResponseUserResponse "Successfully retrieved profile"
// @Failure 404 {object} model.SwaggerWebResponseString "User not found"
// @Security BearerAuth
// @Router /users/profile [get]
func (s *UserServiceImpl) Profile(c *fiber.Ctx) error {
	ctx := c.UserContext()
	val := ctx.Value("user")
	UserId := val.(*model.Claims).UserID

	Users, err := s.repoUser.FindById(ctx, UserId)
	if err != nil {
		response := model.WebResponse[string]{
			Status: "error",
			Data:   "User not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	UserResponse := model.UserResponse{
		ID:       Users.User.ID,
		Username: Users.User.Username,
		Email:    Users.User.Email,
		FullName: Users.User.FullName,
		Role:     Users.User.RoleName,
	}

	if Users.StudentID.Valid {
		UserResponse.StudentProfile = &model.StudentCreate{
			StudentID:    Users.StudentID.String,
			ProgramStudy: Users.ProgramStudy.String,
			AcademicYear: Users.AcademicYear.String,
			AdvisorID:    Users.AdvisorID.String,
		}
	} else if Users.LecturerID.Valid {
		UserResponse.LecturerProfile = &model.LecturerCreate{
			LecturerID: Users.LecturerID.String,
			Department: Users.Department.String,
		}
	}

	response := model.WebResponse[model.UserResponse]{
		Status: "success",
		Data:   UserResponse,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
