package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"go_dfs/internal/model"
	"go_dfs/internal/redis"
	"go_dfs/internal/repository"
	"go_dfs/internal/storage"
	"io"
	"time"

	"github.com/google/uuid"
)

type PublicFileShareService struct{
	fileRepository *repository.FileRepository
	publicFileShareRepository *repository.PublicFileShareRepository
	publicShareCache *redis.PublicShareCache
	storage storage.Storage
}

func NewPublicFileShareService(
	fileRepository *repository.FileRepository,
	publicFileShareRepository *repository.PublicFileShareRepository,
	publicShareCache *redis.PublicShareCache,
	storage storage.Storage,
) *PublicFileShareService{
	return &PublicFileShareService{
		fileRepository:fileRepository ,
		publicFileShareRepository: publicFileShareRepository,
		publicShareCache:publicShareCache ,
		storage:storage,
	}
}

func(s *PublicFileShareService) CreatePublicShare(
	ctx context.Context,
	fileID uuid.UUID,
	ownerID uuid.UUID,
	expiresAt *time.Time,
) (*model.PublicShare,error){
	file,err:=s.fileRepository.GetFileByID(ctx,fileID,ownerID)
	if err!=nil{
		return nil,err
	}
	if file==nil{
		return nil,ErrFileNotFound
	}

	token,err:=generatePublicShareToken()
	if err!=nil{
		return nil,err
	}

	share:=&model.PublicShare{
		FileID:fileID,
		OwnerID:ownerID,
		Token:token,
		ExpiresAt:expiresAt,
	}

	err=s.publicFileShareRepository.CreatePublicShare(ctx,share)
	if err!=nil{
		return nil,err
	}
	return share,nil
}

func generatePublicShareToken() (string,error){
	bytes:=make([]byte,32)

	_,err:=rand.Read(bytes)
	if err!=nil{
		return "",errors.New("failed to generate public share token")
	}
	return base64.RawURLEncoding.EncodeToString(bytes),nil
}

func(s *PublicFileShareService) GetPublicFile(
	ctx context.Context,
	token string,
) (*model.File,io.ReadCloser,error){

	if token==""{
		return nil,nil,ErrFileNotFound
	}

	share,err:=s.publicShareCache.Get(ctx, token)
	if err!=nil{
		share=nil
	}
	cacheMiss:= share==nil

	if cacheMiss{
		share,err=s.publicFileShareRepository.GetPublicShareByToken(
			ctx,token,
		)
		if err!=nil{
			return nil,nil,err
		}
		if share==nil{
			return nil,nil,ErrFileNotFound
		}
	}
	if share.RevokedAt!=nil{
		return nil,nil,ErrFileNotFound
	}
	if share.ExpiresAt !=nil && time.Now().After(*share.ExpiresAt){
		return nil,nil,ErrFileNotFound
	}

	if cacheMiss{
		ttl:=10*time.Minute
		if share.ExpiresAt !=nil{
			ttl = time.Until(*share.ExpiresAt)
			if ttl<=0{
				return nil,nil,ErrFileNotFound
			}
		}
	
		_ = s.publicShareCache.Set(ctx,share,ttl,)
	}

	file,err:=s.fileRepository.GetFileByID(
		ctx,
		share.FileID,
		share.OwnerID,
	)
	if err!=nil{
		return nil,nil,err
	}

	reader,err:=s.storage.Get(
		ctx,
		file.StorageKey,
	)
	if err!=nil{
		return nil,nil,err
	}
	return file,reader,nil
}

func(s *PublicFileShareService) RevokePublicShare(
	ctx context.Context,
	fileID uuid.UUID,
	ownerID uuid.UUID,
)error{
	file,err:=s.fileRepository.GetFileByID(
		ctx,
		fileID,
		ownerID,
	)
	if err!=nil{
		return err
	}
	if file==nil{
		return ErrFileNotFound
	}

	share,err:=s.publicFileShareRepository.RevokePublicShare(
		ctx,
		fileID,
		ownerID,
	)
	if err!=nil{
		return err
	}
	if share==nil{
		return ErrFileNotFound
	}

	if share.Token !=""{
		_=s.publicShareCache.Delete(ctx, share.Token,)
	}
	return nil
}