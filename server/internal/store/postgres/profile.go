package store

import (
	"context"
	"errors"

	"github.com/Raaffs/profileManager/server/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)
type PostgresProfileRepo struct {
	Pool *pgxpool.Pool
}

func (r *PostgresProfileRepo) GetProfileByUserID(ctx context.Context, userID int) (*models.Profile, error) {
	var p models.Profile
	query:=`
		SELECT 
		user_id,
		full_name,
		date_of_birth,
		phone_number,
		address,
		aadhaar_number
		FROM profiles
		WHERE user_id=$1
	`
     if err:=r.Pool.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&p.UserID,
		&p.FullName,
		&p.DateOfBirth,
		&p.PhoneNumber,
		&p.Address,
		&p.AadhaarNumber,
	);err!=nil{
		if errors.Is(err, pgx.ErrNoRows){
			return nil, models.NotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *PostgresProfileRepo) Create(ctx context.Context, profile models.Profile) error {
	query:=`
		INSERT INTO profiles (user_id,full_name,date_of_birth,phone_number,address,aadhaar_number)
		VALUES ($1,$2,$3,$4,$5,$6)
	`
	_,err:=r.Pool.Exec(
		ctx,
		query,
		profile.UserID,
		profile.FullName,
		profile.DateOfBirth,
		profile.PhoneNumber,
		profile.Address,
		profile.AadhaarNumber,
	)
	
	if err!=nil{
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return models.AlreadyExists
		}
		return err
	}
	return nil
}

func (r *PostgresProfileRepo) Update(ctx context.Context, profile models.Profile) error {
	query:=`
		UPDATE profiles
		SET full_name=$1,
		    date_of_birth=$2,
			phone_number=$3,
			address=$4,
			aadhaar_number=$5
		WHERE user_id=$6
	`
	_,err:=r.Pool.Exec(
		ctx,
		query,
		profile.FullName,
		profile.DateOfBirth,
		profile.PhoneNumber,
		profile.Address,
		profile.AadhaarNumber,
		profile.UserID,
	);if err!=nil{
		if errors.Is(err, pgx.ErrNoRows){
			return models.NotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return models.AlreadyExists
		}
		return err
	}
	return nil
}

func (r *PostgresProfileRepo) Delete(ctx context.Context, userID int) error {
	return nil
}
