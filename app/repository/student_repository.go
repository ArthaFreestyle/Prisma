package repository

import (
	"context"
	"database/sql"
	"prisma/app/model"

	"github.com/sirupsen/logrus"
)

type StudentRepository interface {
	Save(ctx context.Context, tx *sql.Tx, Student *model.Student) (*model.Student, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]model.Student, error)
	FindById(ctx context.Context, tx *sql.Tx, id string) (*model.Student, error)
	DeleteById(ctx context.Context, tx *sql.Tx, id string) error
}

type StudentRepositoryImpl struct {
	Log *logrus.Logger
}

func (repo *StudentRepositoryImpl) DeleteById(ctx context.Context, tx *sql.Tx, id string) error {
	SQL := "DELETE FROM students WHERE id = $1;"
	_, err := tx.ExecContext(ctx, SQL, id)
	if err != nil {
		return err
	}
	return nil
}

func NewStudentRepositoryImpl(log *logrus.Logger) StudentRepository {
	return &StudentRepositoryImpl{
		Log: log,
	}
}

func (repo *StudentRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, Student *model.Student) (*model.Student, error) {
	SQL := "INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id) VALUES ($1,$2,$3,$4,$5) returning id"
	err := tx.QueryRowContext(ctx, SQL, Student.UserID, Student.StudentID, Student.ProgramStudy, Student.AcademicYear, Student.AdvisorID).Scan(&Student.ID)
	if err != nil {
		return nil, err
	}

	return Student, nil
}

func (repo *StudentRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]model.Student, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *StudentRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (*model.Student, error) {
	//TODO implement me
	panic("implement me")
}
