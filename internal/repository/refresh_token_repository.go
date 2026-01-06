package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(
	ctx context.Context,
	userID string,
	token string,
	expiresAt time.Time,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
	`, uuid.New(), userID, token, expiresAt)

	return err
}

func (r *RefreshTokenRepository) FindActive(
	ctx context.Context,
	token string,
) (string, error) {
	var userID string

	err := r.db.QueryRow(ctx, `
		SELECT user_id
		FROM refresh_tokens
		WHERE token = $1
		  AND status = 'active'
		  AND expires_at > NOW()
	`, token).Scan(&userID)

	return userID, err
}

func (r *RefreshTokenRepository) Revoke(
	ctx context.Context,
	token string,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens
		SET status = 'revoked'
		WHERE token = $1
	`, token)

	return err
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, userID, token string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens
		SET status = 'revoked'
		WHERE user_id = $1 AND token = $2
	`, userID, token)
	return err
}

func (r *RefreshTokenRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens
		SET status = 'revoked'
		WHERE user_id = $1
	`, userID)
	return err
}