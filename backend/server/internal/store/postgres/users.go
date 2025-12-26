package store

import (
	"context"
	"errors"

	"github.com/Raaffs/profileManager/server/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepo struct {
	Pool *pgxpool.Pool
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	query := `
		SELECT id,email,username
		FROM users
	`
	if err := r.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query := `
		SELECT id,email,username,password_hash
		FROM users
		WHERE email=$1
	`
	if err := r.Pool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepo) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email,username,password_hash)
		VALUES ($1,$2,$3)
	`
	_, err := r.Pool.Exec(
		ctx,
		query,
		user.Email,
		user.Username,
		user.PasswordHash,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return models.AlreadyExists
		}
		return err
	}

	return nil
}

