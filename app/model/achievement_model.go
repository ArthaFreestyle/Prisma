package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongo struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StudentID       string             `bson:"studentId" json:"student_id"` // Disimpan sebagai string UUID
	AchievementType string             `bson:"achievementType" json:"achievement_type"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         AchievementDetails `bson:"details" json:"details"`
	Attachments     []Attachment       `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Tags            []string           `bson:"tags" json:"tags"`
	Points          int                `bson:"points,omitempty" json:"points,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updated_at"`
}

type AchievementDetails struct {
	// Competition
	CompetitionName  string `bson:"competitionName,omitempty" json:"competition_name,omitempty"`
	CompetitionLevel string `bson:"competitionLevel,omitempty" json:"competition_level,omitempty"`
	Rank             int    `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType        string `bson:"medalType,omitempty" json:"medal_type,omitempty"`

	// Organization
	OrganizationName string    `bson:"organizationName,omitempty" json:"organization_name,omitempty"`
	Position         string    `bson:"position,omitempty" json:"position,omitempty"`
	StartDate        time.Time `bson:"startDate,omitempty" json:"start_date,omitempty"`
	EndDate          time.Time `bson:"endDate,omitempty" json:"end_date,omitempty"`

	// General
	EventDate time.Time `bson:"eventDate,omitempty" json:"event_date,omitempty"`
	Location  string    `bson:"location,omitempty" json:"location,omitempty"`
	Organizer string    `bson:"organizer,omitempty" json:"organizer,omitempty"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"file_name"`
	FileURL    string    `bson:"fileUrl" json:"file_url"`
	FileType   string    `bson:"fileType" json:"file_type"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploaded_at"`
}

type CreateAchievementRequest struct {
	AchievementType string             `json:"achievement_type" validate:"required"`
	Title           string             `json:"title" validate:"required"`
	Description     string             `json:"description" validate:"required"`
	Details         AchievementDetails `json:"details"`
	Tags            []string           `json:"tags"`
}

type UpdateAchievementRequest struct {
	ID              string             `json:"id" validate:"required"`
	AchievementType string             `json:"achievement_type" validate:"required"`
	Title           string             `json:"title" validate:"required"`
	Description     string             `json:"description" validate:"required"`
	Details         AchievementDetails `json:"details"`
	Tags            []string           `json:"tags"`
}

type AchievementReference struct {
	ID                 string            `json:"id"`
	StudentID          string            `json:"student_id"`
	MongoAchievementID string            `json:"mongo_achievement_id"`
	Status             string            `json:"status"`
	RejectionNote      string            `json:"rejection_note,omitempty"`
	SubmittedAt        *time.Time        `json:"submitted_at,omitempty"`
	VerifiedAt         *time.Time        `json:"verified_at,omitempty"`
	VerifiedBy         string            `json:"verified_by,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	Detail             *AchievementMongo `json:"detail,omitempty"`
}

type AchievementReferenceDetail struct {
	ID                 string            `json:"id"`
	MongoAchievementID string            `json:"mongo_achievement_id"`
	Status             string            `json:"status"`
	RejectionNote      *string           `json:"rejection_note,omitempty"`
	SubmittedAt        *time.Time        `json:"submitted_at,omitempty"`
	VerifiedAt         *time.Time        `json:"verified_at,omitempty"`
	VerifiedBy         *string           `json:"verified_by,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
	Detail             *AchievementMongo `json:"detail,omitempty"`
	UserDetail         UserResponse      `json:"user_detail"`
}

type AchievementReferenceLecturer struct {
	ID                 string            `json:"id"`
	MongoAchievementID string            `json:"-"`
	Student            UserResponse      `json:"student"`
	Title              string            `json:"title"`
	Type               string            `json:"type"`
	Detail             *AchievementMongo `json:"detail,omitempty"`
	Status             string            `json:"status"`
	CreatedAt          time.Time         `json:"created_at"`
}

type AchievementReferenceStudent struct {
	ID                 string            `json:"id"`
	Title              string            `json:"title"`
	Type               string            `json:"type"`
	CreatedAt          time.Time         `json:"created_at"`
	MongoAchievementID string            `json:"-"`
	Detail             *AchievementMongo `json:"detail,omitempty"`
	Status             string            `json:"status"`
}

type AchievementReferenceAdmin struct {
	ID                 string            `json:"id"`
	MongoAchievementID string            `json:"-"`
	Title              string            `json:"title"`
	Type               string            `json:"type"`
	Student            UserResponse      `json:"student"`
	Lecturer           UserResponse      `json:"lecturer"`
	Detail             *AchievementMongo `json:"detail,omitempty"`
	Status             string            `json:"status"`
	CreatedAt          time.Time         `json:"created_at"`
}

type AchievementHistory struct {
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

type CreateRejection struct {
	RejectionNote string `json:"rejection_note"`
}
