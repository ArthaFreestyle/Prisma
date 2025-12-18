package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"prisma/app/model"
	"strings"
	"time"
)

type AchievementReferenceRepository interface {
	Create(ctx context.Context, achievement model.AchievementReference) (*model.AchievementReference, error)
	Update(ctx context.Context, achievement model.AchievementReference) (*model.AchievementReference, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.AchievementReferenceDetail, error)
	FindByLecturer(ctx context.Context, id string, page int, limit int) ([]model.AchievementReferenceLecturer, error)
	FindByStudent(ctx context.Context, id string, page int, limit int) ([]model.AchievementReferenceStudent, error)
	FindAll(ctx context.Context, page int, limit int) ([]model.AchievementReferenceAdmin, error)
}

type achievementReferenceRepository struct {
	Log *logrus.Logger
	DB  *sql.DB
}

func NewAchievementReferenceRepository(log *logrus.Logger, db *sql.DB) AchievementReferenceRepository {
	return &achievementReferenceRepository{
		Log: log,
		DB:  db,
	}
}

func (repo *achievementReferenceRepository) Create(ctx context.Context, achievement model.AchievementReference) (*model.AchievementReference, error) {
	ts := time.Now()
	SQL := "INSERT INTO achievement_references(student_id, mongo_achievement_id, status, created_at,updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := repo.DB.QueryRowContext(ctx, SQL, achievement.StudentID, achievement.MongoAchievementID, achievement.Status, ts, ts).Scan(&achievement.ID)
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

func (repo *achievementReferenceRepository) Update(ctx context.Context, achievement model.AchievementReference) (*model.AchievementReference, error) {
	ts := time.Now()

	var setClauses []string
	args := []interface{}{}
	argId := 1

	if achievement.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argId))
		args = append(args, achievement.Status)
		argId++
	}
	if achievement.SubmittedAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("submitted_at = $%d", argId))
		args = append(args, achievement.SubmittedAt)
		argId++
	}

	if achievement.VerifiedBy != "" {
		setClauses = append(setClauses, fmt.Sprintf("verified_at = $%d", argId))
		args = append(args, ts)
		argId++

		setClauses = append(setClauses, fmt.Sprintf("verified_by = $%d", argId))
		args = append(args, achievement.VerifiedBy)
		argId++
	}

	if achievement.RejectionNote != "" {
		setClauses = append(setClauses, fmt.Sprintf("rejection_note = $%d", argId))
		args = append(args, achievement.RejectionNote)
		argId++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argId))
	args = append(args, ts)
	argId++

	setQuery := strings.Join(setClauses, ", ")

	SQL := fmt.Sprintf("UPDATE achievement_references SET %s WHERE id = $%d", setQuery, argId)
	args = append(args, achievement.ID)

	res, err := repo.DB.ExecContext(ctx, SQL, args...)
	if err != nil {
		return nil, err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if aff == 0 {
		return nil, errors.New("no rows affected")
	}

	return &achievement, nil
}
func (repo *achievementReferenceRepository) Delete(ctx context.Context, id string) error {
	SQL := "UPDATE achievement_references SET status = 'DELETED' WHERE id = $1"
	res, err := repo.DB.ExecContext(ctx, SQL, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repo *achievementReferenceRepository) FindByID(ctx context.Context, id string) (*model.AchievementReferenceDetail, error) {
	SQL := `SELECT a.id,a.status,a.mongo_achievement_id,a.submitted_at,a.verified_at,
     a.verified_by,a.rejection_note,a.created_at,a.updated_at,
    u.username,u.full_name,u.email,s.student_id,s.academic_year,s.program_study FROM achievement_references as a
        JOIN students as s ON s.id = a.student_id
        JOIN users as u ON u.id = s.user_id   
           WHERE a.id = $1 AND a.status != 'DELETED'`

	achievement := model.AchievementReferenceDetail{}

	err := repo.DB.QueryRowContext(ctx, SQL, id).Scan(
		&achievement.ID,
		&achievement.Status,
		&achievement.MongoAchievementID,
		&achievement.SubmittedAt,
		&achievement.VerifiedAt,
		&achievement.VerifiedBy,
		&achievement.RejectionNote,
		&achievement.CreatedAt,
		&achievement.UpdatedAt,
		&achievement.UserDetail.Username,
		&achievement.UserDetail.FullName,
		&achievement.UserDetail.Email,
		&achievement.UserDetail.StudentProfile.StudentID,
		&achievement.UserDetail.StudentProfile.AcademicYear,
		&achievement.UserDetail.StudentProfile.ProgramStudy,
	)

	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

func (repo *achievementReferenceRepository) FindByLecturer(ctx context.Context, id string, page int, limit int) ([]model.AchievementReferenceLecturer, error) {
	skip := (page - 1) * limit
	SQL := `SELECT a.id,a.mongo_achievement_id,a.status,u.username,u.full_name,u.email,
			s.program_study,s.academic_year,s.student_id FROM achievement_references a
			JOIN students as s ON s.id = a.student_id
			JOIN lecturers as l ON l.id = s.advisor_id
            JOIN users as u ON u.id = s.user_id
			WHERE l.user_id = $1 AND a.status != 'DELETED'
			LIMIT $2 OFFSET $3`

	rows, err := repo.DB.QueryContext(ctx, SQL, id, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	achievements := []model.AchievementReferenceLecturer{}
	for rows.Next() {
		achievement := model.AchievementReferenceLecturer{}
		err := rows.Scan(&achievement.ID, &achievement.MongoAchievementID, &achievement.Status,
			&achievement.Student.Username, &achievement.Student.FullName, &achievement.Student.Email,
			&achievement.Student.StudentProfile.ProgramStudy, &achievement.Student.StudentProfile.AcademicYear,
			&achievement.Student.StudentProfile.StudentID)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, achievement)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (repo *achievementReferenceRepository) FindByStudent(ctx context.Context, id string, page int, limit int) ([]model.AchievementReferenceStudent, error) {
	skip := (page - 1) * limit
	SQL := `SELECT a.id,a.mongo_achievement_id,a.status
			FROM achievement_references a
			JOIN students as s ON s.id = a.student_id
            JOIN users as u ON u.id = s.user_id
			WHERE s.user_id = $1 AND COALESCE(a.status, '') != 'DELETED'
			LIMIT $2 OFFSET $3`

	rows, err := repo.DB.QueryContext(ctx, SQL, id, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	achievements := []model.AchievementReferenceStudent{}
	for rows.Next() {
		achievement := model.AchievementReferenceStudent{}
		err := rows.Scan(&achievement.ID, &achievement.MongoAchievementID, &achievement.Status)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, achievement)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (repo *achievementReferenceRepository) FindAll(ctx context.Context, page int, limit int) ([]model.AchievementReferenceAdmin, error) {
	skip := (page - 1) * limit
	SQL := `SELECT a.id,a.mongo_achievement_id,a.status,u.username,u.full_name,u.email,
			s.program_study,s.academic_year,s.student_id,l.department,u2.username,u2.email,u2.full_name FROM achievement_references a
			JOIN students as s ON s.id = a.student_id
			JOIN lecturers as l ON l.id = s.advisor_id
            JOIN users as u ON u.id = s.user_id
			JOIN users as u2 ON u2.id = l.user_id
			WHERE a.status != 'DELETED'                                                                                   
			LIMIT $1 OFFSET $2`

	rows, err := repo.DB.QueryContext(ctx, SQL, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	achievements := []model.AchievementReferenceAdmin{}
	for rows.Next() {
		achievement := model.AchievementReferenceAdmin{}
		err := rows.Scan(&achievement.ID, &achievement.MongoAchievementID, &achievement.Status,
			&achievement.Student.Username, &achievement.Student.FullName, &achievement.Student.Email,
			&achievement.Student.StudentProfile.ProgramStudy, &achievement.Student.StudentProfile.AcademicYear,
			&achievement.Student.StudentProfile.StudentID, &achievement.Lecturer.LecturerProfile.Department, &achievement.Lecturer.Username,
			&achievement.Lecturer.Email, &achievement.Lecturer.FullName)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, achievement)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return achievements, nil
}
