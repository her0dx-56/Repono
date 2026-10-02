package repository

import (
	"context"
	"errors"
	"go_dfs/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PublicFileShareRepository struct{
	db *pgxpool.Pool
}

func NewPublicFileShareRepository(db *pgxpool.Pool) *PublicFileShareRepository{
	return &PublicFileShareRepository{
		db:db,
	}
}

func(r *PublicFileShareRepository) CreatePublicShare(
	ctx context.Context,
	share *model.PublicShare,
) error{
	query:=`
		INSERT INTO public_file_shares(
			file_id,
			owner_id,
			token,
			expires_at
		)
		VALUES ($1,$2,$3,$4)
		RETURNING id,created_at
	`
	err:= r.db.QueryRow(
		ctx,
		query,
		share.FileID,
		share.OwnerID,
		share.Token,
		share.ExpiresAt,
	).Scan(
		&share.ID,
		&share.CreatedAt,
	)
	if err !=nil{
		return err
	}
	return nil
}

func(r *PublicFileShareRepository) GetPublicShareByToken(
	ctx context.Context,
	token string,
)(*model.PublicShare,error){
	query:=`
		SELECT
			id,
			file_id,
			owner_id,
			token,
			expires_at,
			created_at,
			revoked_at
			FROM public_file_shares
			WHERE token =$1
	`

	share:=&model.PublicShare{}

	err:= r.db.QueryRow(
		ctx,
		query,
		token,
	).Scan(
		&share.ID,
		&share.FileID,
		&share.OwnerID,
		&share.Token,
		&share.ExpiresAt,
		&share.CreatedAt,
		&share.RevokedAt,
	)
	if err!=nil{
		if errors.Is(err,pgx.ErrNoRows){
			return nil,nil
		}
		return nil,err
	}
	return share,nil
}

func (r *PublicFileShareRepository) RevokePublicShare(
	ctx context.Context,
	fileID uuid.UUID,
	ownerID uuid.UUID,
)(*model.PublicShare,error){
	query:=`
		UPDATE public_file_shares
		SET revoked_at=NOW()
		WHERE file_id=$1
		AND owner_id=$2
		AND revoked_at IS NUll
		 RETURNING
            id,
            file_id,
            owner_id,
            token,
            expires_at,
            created_at,
            revoked_at
	`
	var share model.PublicShare

    err := r.db.QueryRow(
        ctx,
        query,
        fileID,
        ownerID,
    ).Scan(
        &share.ID,
        &share.FileID,
        &share.OwnerID,
        &share.Token,
        &share.ExpiresAt,
        &share.CreatedAt,
        &share.RevokedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }

        return nil, err
    }

    return &share, nil
}