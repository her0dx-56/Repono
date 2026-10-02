package repository

import (
	"context"
	"errors"
	"go_dfs/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)
type FileRepository struct{
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository{
	return &FileRepository{
		db:db,
	}
}

func (r *FileRepository) CreateFile(
	ctx context.Context,
	file *model.File,
)error{
	query:=`INSERT INTO files(
		user_id,
		name,
		storage_key,
		size,
		content_type
	)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id,created_at,updated_at
	`
	return r.db.QueryRow(
		ctx,
		query,
		file.UserID,
		file.Name,
		file.StorageKey,
		file.Size,
		file.ContentType,
	).Scan(
		&file.ID,
		&file.CreatedAt,
		&file.UpdatedAt,
	)
}

func (r *FileRepository) GetFileByID(
    ctx context.Context,
    fileID uuid.UUID,
    userID uuid.UUID,
) (*model.File, error) {
	query:=`
	  SELECT
	  	 id,
		 user_id,
		 name,
		 storage_key,
		 size,
		 content_type,
		 created_at,
		 updated_at
	   FROM files 
	   WHERE id = $1
	   AND user_id = $2
	`
	file:=&model.File{}

	err:=r.db.QueryRow(
		ctx,
		query,
		fileID,
		userID,
	).Scan(
		&file.ID,
		&file.UserID,
		&file.Name,
		&file.StorageKey,
		&file.Size,
		&file.ContentType,
		&file.CreatedAt,
		&file.UpdatedAt,
	)
	if err !=nil{
		if errors.Is(err,pgx.ErrNoRows){
		 return nil,nil
		}
		return nil,err
	}
	return file,nil
}

func (r *FileRepository) GetFilesByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.File,error){
	query:=`
	   SELECT
	  	 id,
		 user_id,
		 name,
		 storage_key,
		 size,
		 content_type,
		 created_at,
		 updated_at
	   FROM files 
	   WHERE user_id = $1
	   ORDER BY created_at DESC
	`
  files:=[]model.File{}

  rows,err:=r.db.Query(
		ctx,
		query,
		userID,
	)
	if err !=nil{
		return nil,err
	}
	defer rows.Close()

	for rows.Next(){
		file:=model.File{}
		if err:=rows.Scan(
			&file.ID,
			&file.UserID,
			&file.Name,
			&file.StorageKey,
			&file.Size,
			&file.ContentType,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err !=nil{
			return nil,err
		}
		files=append(files,file)
	}

	if err:=rows.Err(); err!=nil{
		return nil,err
	}
	return files,nil
} 

func(r *FileRepository) DeleteFile(
	ctx context.Context,
	fileID uuid.UUID,
	userID uuid.UUID,
)error{
	query:=`
		DELETE FROM files
		WHERE id=$1
		And user_id=$2
	`
	_,err:=r.db.Exec(
		ctx,
		query,
		fileID,
		userID,
	)
	if err !=nil{
		return err
	}
	return nil
}