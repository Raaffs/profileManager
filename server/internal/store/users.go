package store

import (
	"context"

	"github.com/Raaffs/profileManager/server/internal/models"
	"github.com/Raaffs/profileManager/server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepo struct {
	Pool *pgxpool.Pool
}

func (r *PostgresUserRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return &models.User{},nil
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return nil // Logic here
}

func (r *PostgresUserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	return nil // Logic here
}

func (r *PostgresUserRepo) DeleteUser(ctx context.Context, id int) error {
	return nil // Logic here
}

type PostgresProfileRepo struct {
	Pool *pgxpool.Pool
}

func (r *PostgresProfileRepo) GetProfileByUserID(ctx context.Context, userID int) (*models.Profile, error) {
	return &models.Profile{}, nil
}

func (r *PostgresProfileRepo) CreateProfile(ctx context.Context, profile *models.Profile) error {
	return nil
}

func (r *PostgresProfileRepo) UpdateProfile(ctx context.Context, profile *models.Profile) error {
	return nil
}

func (r *PostgresProfileRepo) DeleteProfile(ctx context.Context, userID int) error {
	return nil
}

// Constructor
func NewPostgresRepo(pool *pgxpool.Pool) *repository.Repository {
	return &repository.Repository{
		Users:    &PostgresUserRepo{Pool: pool},
		Profiles: &PostgresProfileRepo{Pool: pool},
	}
}