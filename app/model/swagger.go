package model

// SwaggerUserProfile is a swagger-compatible version of UserProfile
// Used for API documentation only, avoiding sql.NullString types
type SwaggerUserProfile struct {
	ID           string  `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	FullName     string  `json:"fullname"`
	RoleName     string  `json:"role_name"`
	StudentID    *string `json:"student_id,omitempty"`
	ProgramStudy *string `json:"program_study,omitempty"`
	AcademicYear *string `json:"academic_year,omitempty"`
	AdvisorID    *string `json:"advisor_id,omitempty"`
	LecturerID   *string `json:"lecturer_id,omitempty"`
	Department   *string `json:"department,omitempty"`
}

// SwaggerStudent for swagger documentation
type SwaggerStudent struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	StudentID    string  `json:"student_id"`
	ProgramStudy string  `json:"program_study"`
	AcademicYear string  `json:"academic_year"`
	AdvisorID    *string `json:"advisor_id,omitempty"`
}

// SwaggerWebResponseUserProfile for single UserProfile
type SwaggerWebResponseUserProfile struct {
	Status string             `json:"status"`
	Data   SwaggerUserProfile `json:"data"`
	Errors string             `json:"errors,omitempty"`
}

// SwaggerWebResponseUserProfiles for array of UserProfile
type SwaggerWebResponseUserProfiles struct {
	Status string               `json:"status"`
	Data   []SwaggerUserProfile `json:"data"`
	Errors string               `json:"errors,omitempty"`
}

// SwaggerWebResponseStudent for swagger documentation
type SwaggerWebResponseStudent struct {
	Status string         `json:"status"`
	Data   SwaggerStudent `json:"data"`
	Errors string         `json:"errors,omitempty"`
}

// SwaggerWebResponseString for swagger documentation
type SwaggerWebResponseString struct {
	Status string `json:"status"`
	Data   string `json:"data,omitempty"`
	Errors string `json:"errors,omitempty"`
}

// SwaggerWebResponseInterface for swagger documentation
type SwaggerWebResponseInterface struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Errors string      `json:"errors,omitempty"`
}

// SwaggerWebResponseUserResponse for swagger documentation
type SwaggerWebResponseUserResponse struct {
	Status string       `json:"status"`
	Data   UserResponse `json:"data"`
	Errors string       `json:"errors,omitempty"`
}

// SwaggerWebResponseUserResponses for array of UserResponse
type SwaggerWebResponseUserResponses struct {
	Status string         `json:"status"`
	Data   []UserResponse `json:"data"`
	Errors string         `json:"errors,omitempty"`
}

// SwaggerWebResponseUserUpdateResponse for swagger documentation
type SwaggerWebResponseUserUpdateResponse struct {
	Status string             `json:"status"`
	Data   UserUpdateResponse `json:"data"`
	Errors string             `json:"errors,omitempty"`
}

// SwaggerWebResponseAchievementReferenceAdmin for swagger documentation
type SwaggerWebResponseAchievementReferenceAdmin struct {
	Status string                      `json:"status"`
	Data   []AchievementReferenceAdmin `json:"data"`
	Errors string                      `json:"errors,omitempty"`
}

// SwaggerChangeAdvisorRequest for swagger documentation
type SwaggerChangeAdvisorRequest struct {
	AdvisorID string `json:"advisor" example:"uuid-of-advisor"`
}

// ConvertToSwaggerUserProfile converts UserProfile to SwaggerUserProfile
func ConvertToSwaggerUserProfile(up *UserProfile) *SwaggerUserProfile {
	swagger := &SwaggerUserProfile{
		ID:       up.User.ID,
		Username: up.User.Username,
		Email:    up.User.Email,
		FullName: up.User.FullName,
		RoleName: up.User.RoleName,
	}

	if up.StudentID.Valid {
		swagger.StudentID = &up.StudentID.String
	}
	if up.ProgramStudy.Valid {
		swagger.ProgramStudy = &up.ProgramStudy.String
	}
	if up.AcademicYear.Valid {
		swagger.AcademicYear = &up.AcademicYear.String
	}
	if up.AdvisorID.Valid {
		swagger.AdvisorID = &up.AdvisorID.String
	}
	if up.LecturerID.Valid {
		swagger.LecturerID = &up.LecturerID.String
	}
	if up.Department.Valid {
		swagger.Department = &up.Department.String
	}

	return swagger
}

// ConvertToSwaggerStudent converts Student to SwaggerStudent
func ConvertToSwaggerStudent(s *Student) *SwaggerStudent {
	swagger := &SwaggerStudent{
		ID:           s.ID,
		UserID:       s.UserID,
		StudentID:    s.StudentID,
		ProgramStudy: s.ProgramStudy,
		AcademicYear: s.AcademicYear,
	}

	if s.AdvisorID != "" {
		swagger.AdvisorID = &s.AdvisorID
	}

	return swagger
}
