package repository

import (
	"context"
	"database/sql"
	"personal_expense_tracker/domain"
)

type UserRepository interface {
	Add(ctx context.Context, user *domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Add(ctx context.Context, user *domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`
	err := u.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1;
	`
	var user domain.User
	err := u.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
