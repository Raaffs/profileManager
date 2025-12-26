package store

import (
	"github.com/Raaffs/profileManager/server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresRepo(pool *pgxpool.Pool) *repository.Repository {
	return &repository.Repository{
		Users:    &PostgresUserRepo{Pool: pool},
		Profiles: &PostgresProfileRepo{Pool: pool},
	}
}