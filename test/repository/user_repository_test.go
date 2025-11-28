package repository

import (
	"context"

	"prisma/app/model"
	"prisma/app/repository"
	"prisma/config"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var Config = config.NewViper()
var Log = config.NewLog(Config)
var DB = config.PostgresConnect(Config, Log)

func TestUserRepository_FindByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.UserRepositoryImpl{DB: db}

	ctx := context.Background()
	username := "alan"

	expectedUser := model.User{
		ID:           "11111111-1111-1111-1111-111111111111",
		Username:     "alan",
		FullName:     "Alan Pratama",
		PasswordHash: "$2y$10$abc",
		RoleName:     "student",
	}

	// JSON hasil JSON_AGG
	permJSON := `[{"resource":"users","action":"read"}]`

	rows := sqlmock.NewRows([]string{
		"id", "username", "full_name", "password_hash", "name", "permissions",
	}).AddRow(
		expectedUser.ID,
		expectedUser.Username,
		expectedUser.FullName,
		expectedUser.PasswordHash,
		expectedUser.RoleName,
		permJSON,
	)

	// Cocokin panggilan SELECT (cukup prefix-nya aja)
	mock.ExpectQuery(`SELECT u.id`).
		WithArgs(username).
		WillReturnRows(rows)

	user, err := repo.FindByUsername(ctx, username)

	require.NoError(t, err)
	require.NotNil(t, user)

	require.Equal(t, expectedUser.ID, user.ID)
	require.Equal(t, expectedUser.Username, user.Username)
	require.Equal(t, expectedUser.FullName, user.FullName)
	require.Equal(t, expectedUser.PasswordHash, user.PasswordHash)
	require.Equal(t, expectedUser.RoleName, user.RoleName)

	//require.Len(t, user.Permissions, 1)
	//require.Equal(t, "users", user.Permissions[0].Resource)
	//require.Equal(t, "read", user.Permissions[0].Action)

	require.NoError(t, mock.ExpectationsWereMet())
}
