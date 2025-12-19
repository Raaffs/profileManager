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
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type ProfileRepository interface {
	GetProfileByUserID(ctx context.Context, userID int) (*models.Profile, error)
	CreateProfile(ctx context.Context, profile models.Profile) error
	UpdateProfile(ctx context.Context, profile models.Profile) error
	DeleteProfile(ctx context.Context, userID int) error
}