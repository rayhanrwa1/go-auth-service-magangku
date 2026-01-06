package repository

import (
	"context"
	"errors"
	"strings"

	"auth-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserExists = errors.New("user already exists")

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id,
			name,
			email,
			password,
			userable_type,
			terms_accepted
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.DB.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.UserableType,
		user.TermsAccepted,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return ErrUserExists
		}
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {

	query := `
		SELECT
			id,
			name,
			email,
			password,
			userable_type,
			terms_accepted,
			created_at
		FROM users
		WHERE email = $1
	`

	row := r.DB.QueryRow(ctx, query, email)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.UserableType,
		&user.TermsAccepted,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2
	`
	_, err := r.DB.Exec(ctx, query, hashedPassword, userID)
	return err
}