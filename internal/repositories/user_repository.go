package repositories

import (
	"context"
	"database/sql"
	"fluxera/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `INSERT INTO users (email , password_hash) VALUES ($1,$2) RETURNING id, created_at`,
		user.Email,
		user.PasswordHash,
	)
	if err := row.Scan(&user.ID, &user.CreatedAt); err != nil {
		return nil, err
	}

	return user, nil

}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := models.User{}

	row := r.db.QueryRowContext(ctx, `SELECT id, email , password_hash, created_at FROM users WHERE email  = $1`, email)

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id int64) (*models.User, error) {
	user := models.User{}

	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`, id)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
