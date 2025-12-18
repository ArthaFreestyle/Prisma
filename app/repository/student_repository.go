package repository

import (
	"context"
	"database/sql"
	"prisma/app/model"

	"github.com/sirupsen/logrus"
)

type StudentRepository interface {
	Save(ctx context.Context, tx *sql.Tx, Student *model.Student) (*model.Student, error)
	FindAll(ctx context.Context) ([]model.UserProfile, error)
	FindById(ctx context.Context, id string) (*model.UserProfile, error)
	FindByUserId(ctx context.Context, userid string) (*model.Student, error)
	DeleteById(ctx context.Context, tx *sql.Tx, id string) error
	UpdateById(ctx context.Context, Student *model.Student) (*model.Student, error)
}

type StudentRepositoryImpl struct {
	Log *logrus.Logger
	DB  *sql.DB
}

func NewStudentRepositoryImpl(log *logrus.Logger, DB *sql.DB) StudentRepository {
	return &StudentRepositoryImpl{
		Log: log,
		DB:  DB,
	}
}

func (repo *StudentRepositoryImpl) UpdateById(ctx context.Context, Student *model.Student) (*model.Student, error) {
	SQL := `UPDATE students SET advisor_id = $1 WHERE id = $2`
	res, err := repo.DB.ExecContext(ctx, SQL, Student.AdvisorID, Student.ID)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, nil
	}
	return Student, nil
}

func (repo *StudentRepositoryImpl) DeleteById(ctx context.Context, tx *sql.Tx, id string) error {
	SQL := "DELETE FROM students WHERE id = $1;"
	_, err := tx.ExecContext(ctx, SQL, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *StudentRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, Student *model.Student) (*model.Student, error) {
	SQL := "INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id) VALUES ($1,$2,$3,$4,$5) returning id"
	err := tx.QueryRowContext(ctx, SQL, Student.UserID, Student.StudentID, Student.ProgramStudy, Student.AcademicYear, Student.AdvisorID).Scan(&Student.ID)
	if err != nil {
		return nil, err
	}

	return Student, nil
}

func (repo *StudentRepositoryImpl) FindAll(ctx context.Context) ([]model.UserProfile, error) {
	SQL := `SELECT u.username,u.email,u.full_name,s.id,s.student_id,s.program_study,s.academic_year,s.advisor_id
	FROM students s
	JOIN users u ON u.id = s.user_id`
	rows, err := repo.DB.QueryContext(ctx, SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.UserProfile
	for rows.Next() {
		var Student model.UserProfile
		rows.Scan(
			&Student.User.Username, &Student.User.Email, &Student.User.FullName,
			&Student.User.ID, &Student.StudentID, &Student.ProgramStudy,
			&Student.AcademicYear, &Student.AdvisorID)
		users = append(users, Student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *StudentRepositoryImpl) FindByUserId(ctx context.Context, id string) (*model.Student, error) {
	Student := model.Student{}
	SQL := "SELECT id,student_id,program_study,academic_year,advisor_id,created_at  FROM students WHERE user_id = $1;"
	err := repo.DB.QueryRowContext(ctx, SQL, id).Scan(&Student.ID, &Student.StudentID, &Student.ProgramStudy, &Student.AcademicYear, &Student.AdvisorID, &Student.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &Student, nil
}

func (repo *StudentRepositoryImpl) FindById(ctx context.Context, id string) (*model.UserProfile, error) {
	Student := model.UserProfile{}
	SQL := `SELECT u.username,u.email,u.full_name,s.id,s.student_id,s.program_study,s.academic_year,s.advisor_id
			FROM students s
			JOIN users u ON u.id = s.user_id
			WHERE s.id = $1`
	err := repo.DB.QueryRowContext(ctx, SQL, id).Scan(&Student.User.Username, &Student.User.Email, &Student.User.FullName,
		&Student.User.ID, &Student.StudentID, &Student.ProgramStudy,
		&Student.AcademicYear, &Student.AdvisorID)
	if err != nil {
		return nil, err
	}
	return &Student, nil
}
