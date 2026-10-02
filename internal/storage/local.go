package storage

import (
	"os"
	"path/filepath"
	"io"
	"fmt"
	"context"
	"errors"
	"strings"
)
type LocalStorage struct{
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage{
	return &LocalStorage{
		basePath:basePath,
	}
}

func(s *LocalStorage) path(key string) (string,error){
	basePath,err:=filepath.Abs(s.basePath)
	if err!=nil{
		return "",err
	}

	path,err:=filepath.Abs(filepath.Join(basePath,key))
	if err!=nil{
		return "",err
	}
	
	relativePath,err:=filepath.Rel(basePath,path)
	if err!=nil{
		return "",err
	}

	if relativePath==".."||strings.HasPrefix(relativePath,".."+string(os.PathSeparator)){
		return "",errors.New("invalid storage key")
	}
	return path,nil
}


func(s *LocalStorage) Save(
	ctx context.Context,
	key string,
	r io.Reader,
)error{
	path,err:=s.path(key)
	if err!=nil{
		return err
	}
	err=os.MkdirAll(filepath.Dir(path),0755)
	if err !=nil{
		return err
	}
	file,err:=os.Create(path)
	if err!=nil{
		return err
	}
	defer file.Close()

	_,err=io.Copy(file,r)
	if err!=nil{
		return err
	}
	return nil
}

func(s *LocalStorage) Get(
	ctx context.Context,
	key string,
) (io.ReadCloser,error){
	path,err:=s.path(key)
	if err!=nil{
		return nil, err
	}
	file,err:=os.Open(path)
	if err!=nil{
		return nil,err
	}
	return file,nil
}

func(s *LocalStorage) Delete(
	ctx context.Context,
	key string,
)error{
	path,err:=s.path(key)
	if err!=nil{
		return err
	}

	err=os.Remove(path)
	if err!=nil{
		return fmt.Errorf("failed to delete file:%w",err)
	}
	return nil
}