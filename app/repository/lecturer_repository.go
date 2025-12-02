package repository

import (
	"context"
	"database/sql"
	"prisma/app/model"
)

type LecturerRepository interface {
	Save(ctx context.Context, tx *sql.Tx, Lecturer *model.Lecturer) (*model.Lecturer, error)
}
