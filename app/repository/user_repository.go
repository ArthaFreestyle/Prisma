package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"prisma/app/model"
)

type UserRepository interface {
	Create(ctx context.Context, User model.UserCreateRequest) (int64, error)
	Update(ctx context.Context, User model.UserUpdateRequest) error
	Delete(ctx context.Context, UserId int) error
	FindById(ctx context.Context, UserId int64) (*model.UserResponse, error)
	FindAll(ctx context.Context) (*[]model.UserResponse, error)
	FindByEmail(ctx context.Context, Email string) (model.User, error)
	Logout()
}

type UserRepositoryImpl struct {
	DB *sql.DB
}

func (repo UserRepositoryImpl) Create(ctx context.Context, User model.UserCreateRequest) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (repo UserRepositoryImpl) Update(ctx context.Context, User model.UserUpdateRequest) error {
	//TODO implement me
	panic("implement me")
}

func (repo UserRepositoryImpl) Delete(ctx context.Context, UserId int) error {
	//TODO implement me
	panic("implement me")
}

func (repo UserRepositoryImpl) FindById(ctx context.Context, UserId int64) (*model.UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (repo UserRepositoryImpl) FindAll(ctx context.Context) (*[]model.UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (repo UserRepositoryImpl) FindByEmail(ctx context.Context, Email string) (*model.User, error) {
	SQL := `SELECT u.id,u.username,u.full_name,u.password_hash,r.name,
			COALESCE(
					JSON_AGG(
						JSON_BUILD_OBJECT('resource', p.resource, 'action', p.action)
					) FILTER (WHERE p.id IS NOT NULL), 
					'[]'
				) as permissions
			FROM users u 
    		INNER JOIN roles r ON u.role_id = r.id
			LEFT JOIN role_permissions rp ON u.role_id = rp.role_id
			LEFT JOIN permissions p ON rp.permission_id = p.id
			WHERE u.email = ?
			GROUP BY u.id,u.username,u.full_name,u.password_hash,r.name;`

	rows, err := repo.DB.Query(SQL, Email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user model.User
	var permissions []byte
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.FullName,
			&user.PasswordHash,
			&user.Role,
			&user.Permissions,
		)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(permissions, &user.Permissions); err != nil {
			return nil, fmt.Errorf("unmarshal permissions: %w", err)
		}
	}

	return &user, nil
}

func (repo UserRepositoryImpl) Logout() {
	//TODO implement me
	panic("implement me")
}
