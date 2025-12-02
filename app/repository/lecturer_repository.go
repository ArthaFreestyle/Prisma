package repository

import (
	"context"
	"database/sql"
	"prisma/app/model"

	"github.com/sirupsen/logrus"
)

type LecturerRepository interface {
	Save(ctx context.Context, tx *sql.Tx, Lecturer *model.Lecturer) (*model.Lecturer, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]model.Lecturer, error)
	FindById(ctx context.Context, tx *sql.Tx, id string) (*model.Lecturer, error)
}

type LecturerRepositoryImpl struct {
	Log *logrus.Logger
}

func (repo *LecturerRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, Lecturer *model.Lecturer) (*model.Lecturer, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *LecturerRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]model.Lecturer, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *LecturerRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (*model.Lecturer, error) {
	//TODO implement me
	panic("implement me")
}
