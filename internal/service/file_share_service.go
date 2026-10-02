package service

import (
	"context"
	"go_dfs/internal/model"
	"go_dfs/internal/repository"
	"errors"
	"github.com/google/uuid"
	"strings"
)

type FileShareService struct{
	fileRepository *repository.FileRepository
	fileShareRepository *repository.FileShareRepository
	userRepository *repository.UserRepository
}

func NewFileShareService(
	fileRepository *repository.FileRepository,
	fileShareRepository *repository.FileShareRepository,
	userRepository *repository.UserRepository,
) *FileShareService{
	return &FileShareService{
		fileRepository: fileRepository,
		fileShareRepository: fileShareRepository,
		userRepository: userRepository,
	}
}
var ErrUserNotFound = errors.New("user not found")
var ErrAlreadyShared=errors.New("file already shared with user")

func(s *FileShareService) CreateShare(
	ctx context.Context,
	fileID uuid.UUID,
	ownerID uuid.UUID,
	shareWithUserID uuid.UUID,
) (*model.FileShare,error){

	file,err:=s.fileRepository.GetFileByID(
		ctx,
		fileID,
		ownerID,
	)
	if err!=nil{
		return nil,err
	}

	if file==nil{
		return nil,ErrFileNotFound
	}

	recipient,err:=s.userRepository.GetUserByID(
		ctx,
		shareWithUserID,
	)
	if err!=nil{
		return nil,err
	}

	if recipient==nil{
		return nil,ErrUserNotFound
	}

	share:=&model.FileShare{
		FileID: fileID,
		OwnerID: ownerID,
		SharedWithUserID: shareWithUserID,
		Permission: "read",
	}

	err=s.fileShareRepository.CreateShare(ctx,share)
	if err!=nil{
		if errors.Is(err,repository.ErrShareAlreadyExists){
			return nil,ErrAlreadyShared
		}
		return nil,err
	}
	return share,nil
}

func(s *FileShareService) GetSharedFiles(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.SharedFileResponse,error){
	return s.fileShareRepository.GetShareByUser(ctx,userID)
}

func(s *FileShareService) CreateShareByEmail(
	ctx context.Context,
	fileID uuid.UUID,
	ownerID uuid.UUID,
	email string,
) (*model.FileShare, error){

	email=strings.TrimSpace(email)
	if email == ""{
		 return nil, errors.New("recipient email is required")
	}

	recipient,err:=s.userRepository.GetUserByEmail(
		ctx,
		email,
	)
	if err!=nil{
		return nil,err
	}
	if recipient == nil{
		return nil,ErrUserNotFound
	}
	if recipient.ID == ownerID{
		return nil, errors.New("cannot share a file with yourself")
	}

	return s.CreateShare(
		ctx,
		fileID,
		ownerID,
		recipient.ID,
	)
}