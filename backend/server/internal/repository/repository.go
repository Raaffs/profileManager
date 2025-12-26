package repository

import (
	"context"
	"github.com/Raaffs/profileManager/server/internal/models"
)

type Repository struct {
	Users    UserRepository
	Profiles ProfileRepository
}

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

type ProfileRepository interface {
	GetByUserID(ctx context.Context, userID int) (*models.Profile, error)
	Create(ctx context.Context, profile models.Profile) error
	Update(ctx context.Context, profile models.Profile) error
}