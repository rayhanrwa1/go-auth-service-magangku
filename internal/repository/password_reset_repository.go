package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type PasswordResetRepository struct {
	DB *pgxpool.Pool
}

func NewPasswordResetRepository(db *pgxpool.Pool) *PasswordResetRepository {
	return &PasswordResetRepository{DB: db}
}

func (r *PasswordResetRepository) Create(ctx context.Context, userID, token string, exp time.Time) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.DB.Exec(ctx, query, userID, token, exp)
	return err
}

func (r *PasswordResetRepository) Verify(ctx context.Context, token string) (string, error) {
	query := `
		SELECT user_id, expires_at
		FROM password_resets
		WHERE token = $1
	`

	var userID string
	var expiresAt time.Time

	err := r.DB.QueryRow(ctx, query, token).Scan(&userID, &expiresAt)
	if err != nil {
		return "", ErrInvalidToken
	}

	if time.Now().After(expiresAt) {
		return "", ErrInvalidToken
	}

	return userID, nil
}

func (r *PasswordResetRepository) Delete(ctx context.Context, token string) error {
	query := `
		DELETE FROM password_resets
		WHERE token = $1
	`
	_, err := r.DB.Exec(ctx, query, token)
	return err
}