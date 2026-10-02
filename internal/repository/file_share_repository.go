package repository

import (
	"context"
	"errors"
	"go_dfs/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileShareRepository struct{
	db *pgxpool.Pool
}

var ErrShareAlreadyExists=errors.New("share already exists")

func NewFileShareRepository (db *pgxpool.Pool) *FileShareRepository{
	return &FileShareRepository{
		db:db,
	}
}

func (r *FileShareRepository) CreateShare(
	ctx context.Context,
	share *model.FileShare,
)error{
	query:=`
		INSERT INTO file_shares(
			file_id,
			owner_id,
			shared_with_user_id,
			permission
		)
		VALUES($1, $2, $3, $4)
		RETURNING id,created_at;
	`
	err:= r.db.QueryRow(
		ctx,
		query,
		share.FileID,
		share.OwnerID,
		share.SharedWithUserID,
		share.Permission,
	).Scan(
		&share.ID,
		&share.CreatedAt,
	)

	if err!=nil{
		var pgErr *pgconn.PgError

		if errors.As(err,&pgErr)&&pgErr.Code=="23505"{
			return ErrShareAlreadyExists
		}
		return err
	}
	return nil
}


func (r *FileShareRepository) GetShare(
	ctx context.Context,
	fileID uuid.UUID,
	userID uuid.UUID,
) (*model.FileShare,error){
	query:=`
		SELECT
			id,
			file_id,
			owner_id,
			shared_with_user_id,
			permission,
			created_at
		FROM file_shares
		WHERE file_id=$1
		AND shared_with_user_id = $2;
	`
	share:= &model.FileShare{}

	err:= r.db.QueryRow(
		ctx,
		query,
		fileID,
		userID,
	).Scan(
		&share.ID,
		&share.FileID,
		&share.OwnerID,
		&share.SharedWithUserID,
		&share.Permission,
		&share.CreatedAt,
	)

	if err!=nil{
		if errors.Is(err,pgx.ErrNoRows){
			return nil,nil
		}
		return nil,err
	}
	return share,nil
}

func (r *FileShareRepository) GetShareByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.SharedFileResponse,error){
	query:=`
		SELECT
			fs.id,
			fs.file_id,
			fs.owner_id,
			fs.Permission ,
			f.name,
			f.size,
			f.content_type,
			f.created_at,
			fs.created_at
		FROM file_shares fs
		JOIN files f
			ON f.id=fs.file_id
		WHERE fs.shared_with_user_id=$1
		ORDER BY fs.created_at DESC;
	`
	rows,err:=r.db.Query(ctx,query,userID)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()

	shares:=make([]model.SharedFileResponse,0)

	for rows.Next(){
		share:=model.SharedFileResponse{}

		err:=rows.Scan(
			&share.ShareID,
			&share.FileID,
			&share.OwnerID,
			&share.Permission,
			&share.Name,
			&share.Size,
			&share.ContentType,
			&share.CreatedAt,
			&share.SharedAt,
		)
		if err!=nil{
			return nil,err
		}
		shares=append(shares,share)
	}

	if err:=rows.Err();err!=nil{
		return nil,err
	}
	return shares,nil
}