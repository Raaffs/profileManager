package store

import (
	"context"

	"github.com/Raaffs/profileManager/server/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)
type PostgresProfileRepo struct {
	Pool *pgxpool.Pool
}

func (r *PostgresProfileRepo) GetProfileByUserID(ctx context.Context, userID int) (*models.Profile, error) {
	return &models.Profile{}, nil
}

func (r *PostgresProfileRepo) CreateProfile(ctx context.Context, profile models.Profile) error {
	return nil
}

func (r *PostgresProfileRepo) UpdateProfile(ctx context.Context, profile models.Profile) error {
	return nil
}

func (r *PostgresProfileRepo) DeleteProfile(ctx context.Context, userID int) error {
	return nil
}
