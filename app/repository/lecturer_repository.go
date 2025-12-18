package repository

import (
	"context"
	"database/sql"
	"prisma/app/model"

	"github.com/sirupsen/logrus"
)

type LecturerRepository interface {
	Save(ctx context.Context, tx *sql.Tx, Lecturer *model.Lecturer) (*model.Lecturer, error)
	FindAll(ctx context.Context) ([]model.UserProfile, error)
	FindById(ctx context.Context, id string) (*model.UserProfile, error)
	DeleteById(ctx context.Context, tx *sql.Tx, id string) error
	FindAllAdvices(ctx context.Context, id string) ([]model.UserProfile, error)
}

type LecturerRepositoryImpl struct {
	Log *logrus.Logger
	DB  *sql.DB
}

func NewLecturerRepositoryImpl(log *logrus.Logger, DB *sql.DB) LecturerRepository {
	return &LecturerRepositoryImpl{
		Log: log,
		DB:  DB,
	}
}

func (repo *LecturerRepositoryImpl) DeleteById(ctx context.Context, tx *sql.Tx, id string) error {
	SQL := "DELETE FROM lecturers WHERE id = $1;"
	_, err := tx.ExecContext(ctx, SQL, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *LecturerRepositoryImpl) FindAllAdvices(ctx context.Context, id string) ([]model.UserProfile, error) {
	SQL := `SELECT u.username,u.email,u.full_name,s.id,s.student_id,s.program_study,s.academic_year,s.advisor_id
	FROM students s
	JOIN users u ON u.id = s.user_id
	WHERE s.advisor_id = $1`

	rows, err := repo.DB.QueryContext(ctx, SQL, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lecturers := []model.UserProfile{}
	for rows.Next() {
		var user model.UserProfile
		rows.Scan(&user.User.Username, &user.User.Email, &user.User.FullName,
			&user.LecturerID, &user.Department)
		lecturers = append(lecturers, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return lecturers, nil
}

func (repo *LecturerRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, Lecturer *model.Lecturer) (*model.Lecturer, error) {
	SQL := "INSERT INTO lecturers (user_id, lecturer_id, department) VALUES ($1,$2,$3) returning id"
	err := tx.QueryRowContext(ctx, SQL, Lecturer.UserID, Lecturer.LecturerID, Lecturer.Department).Scan(&Lecturer.ID)
	if err != nil {
		return nil, err
	}
	return Lecturer, nil
}

func (repo *LecturerRepositoryImpl) FindAll(ctx context.Context) ([]model.UserProfile, error) {
	SQL := `SELECT u.username,u.email,u.full_name,l.id,l.department 
			FROM lecturers l
			JOIN users u ON l.user_id = u.id`

	rows, err := repo.DB.QueryContext(ctx, SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lecturers := []model.UserProfile{}
	for rows.Next() {
		var user model.UserProfile
		rows.Scan(&user.User.Username, &user.User.Email, &user.User.FullName,
			&user.LecturerID, &user.Department)
		lecturers = append(lecturers, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return lecturers, nil
}

func (repo *LecturerRepositoryImpl) FindById(ctx context.Context, id string) (*model.UserProfile, error) {
	SQL := `SELECT u.username,u.email,u.full_name,l.id,l.department 
			FROM lecturers l
			JOIN users u ON l.user_id = u.id
			WHERE l.id = $1`
	row := repo.DB.QueryRowContext(ctx, SQL, id)
	var user model.UserProfile
	row.Scan(&user.User.Username, &user.User.Email, &user.User.FullName, &user.LecturerID, &user.Department)

	return &user, nil
}
