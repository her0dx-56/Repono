package repository

import (
	"context"
	"errors"
	"go_dfs/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct{
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error{
	query:=`INSERT INTO users (
				username,
				email,
				password_hash
			)
			VALUES ($1, $2, $3)
			RETURNING id,created_at, updated_at
	`
	return r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

func(r *UserRepository) GetUserByEmail(ctx context.Context ,email string) (*model.User, error){
	query:= `
		SELECT
			id,
			username,
			email,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE email= $1
	`
	user := &model.User{}

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email, 
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err !=nil{
		if errors.Is(err,pgx.ErrNoRows){
			return nil,nil
		}
		return nil,err
	}

	return user, nil
}

func(r *UserRepository) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (*model.User,error){
	query:=`
		SELECT 
			id,
			username,
			email,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE id=$1
	`
	user:=&model.User{}

	err:=r.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err !=nil{
		if errors.Is(err,pgx.ErrNoRows){
			return nil,nil
		}
		return nil,err
	}
	return user,nil
}

