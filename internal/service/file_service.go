package service

import (
	"context"
	"go_dfs/internal/model"
	"go_dfs/internal/repository"
	"go_dfs/internal/storage"
	"io"
	"errors"
	"github.com/google/uuid"
)

type FileService struct{
	fileRepository *repository.FileRepository
	storage storage.Storage
	fileShareRepository *repository.FileShareRepository
}

func NewFileService (
	fileRepository *repository.FileRepository,
	  fileShareRepository *repository.FileShareRepository,
	storage storage.Storage,
) *FileService{
	return &FileService{
		fileRepository: fileRepository,
		 fileShareRepository :  fileShareRepository ,
		storage:storage,
	}
}

var ErrFileNotFound=errors.New("file not found")

func(s *FileService) UploadFile(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	size int64,
	contentType string,
	r io.Reader,
) (*model.File,error){
	fileID:=uuid.New()
	storageKey:="users/"+ userID.String()+ "/"+fileID.String()

	file:=&model.File{
		ID: fileID,
		UserID: userID,
		StorageKey: storageKey,
		Name: name,
		Size: size,
		ContentType: contentType,
	}

	err:=s.storage.Save(ctx,file.StorageKey,r)
	if err!=nil{
		return nil,err
	}

	err=s.fileRepository.CreateFile(ctx,file)
	if err!=nil{
		_=s.storage.Delete(ctx,file.StorageKey)
		return nil,err
	}
	return file,nil
}

func(s *FileService) GetFile(
	ctx context.Context,
	fileID uuid.UUID,
	userID uuid.UUID,
) (*model.File,io.ReadCloser,error){
	file,err:=s.fileRepository.GetFileByID(
		ctx,
		fileID,
		userID,
	)
	if err!=nil{
		return nil,nil,err
	}
	if file==nil{
		share,err:=s.fileShareRepository.GetShare(
			ctx,
			fileID,
			userID,
		)
		if err!=nil{
			return nil,nil,err
		}
		if share==nil{
			return nil,nil,ErrFileNotFound
		}
		if share.Permission !="read"{
			return nil,nil,ErrFileNotFound 
		}

		file,err=s.fileRepository.GetFileByID(
			ctx,
			fileID,
			share.OwnerID,
		)
		if err!=nil{
			return nil,nil,err
		}

		if file==nil{
			return nil,nil,ErrFileNotFound
		}
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

func(s *FileService) ListFiles(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.File,error){
	files,err:=s.fileRepository.GetFilesByUser(
		ctx,
		userID,
	)
	if err!=nil{
		return nil,err
	}
	return files,err
}

func(s *FileService) DeleteFile(
	ctx context.Context,
	fileID uuid.UUID,
	userID uuid.UUID,
)error{
	file,err:=s.fileRepository.GetFileByID(
		ctx,
		fileID,
		userID,
	)
	if err!=nil{
		return err
	}

	if file==nil{
		return ErrFileNotFound
	}

	err=s.fileRepository.DeleteFile(
		ctx,
		fileID,
		userID,
	)
	if err!=nil{
		return err
	}

	err=s.storage.Delete(
		ctx,
		file.StorageKey,
	)
	if err!=nil{
		return err
	}
	return nil
}
