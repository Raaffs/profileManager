package store

import (
	"context"

	"github.com/Raaffs/profileManager/server/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepo struct {
	Pool *pgxpool.Pool
}

func (r *PostgresUserRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	query:=`
		SELECT id,email,username
		FROM users
	`
	if err:=r.Pool.QueryRow(ctx,query,id).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
	); err!=nil{
		return nil,err
	}
	return &u,nil
}

func (r *PostgresUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query:=`
		SELECT id,email,username
		FROM users
		WHERE email=$1
	`
	if err:=r.Pool.QueryRow(ctx,query,email).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
	); err!=nil{
		return nil,err
	}
	return &u,nil
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	query:=`
		INSERT INTO users (email,username,password_hash)
		VALUES ($1,$2,$3)
	`
	_,err:=r.Pool.Exec(
		ctx,
		query,
		user.Email,
		user.Username,
		user.PasswordHash,
	)
	
	if err!=nil{
		return err
	}

	return nil 
}

func (r *PostgresUserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	return nil 
}

func (r *PostgresUserRepo) DeleteUser(ctx context.Context, id int) error {
	return nil 
}


