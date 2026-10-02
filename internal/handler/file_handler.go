package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_dfs/internal/middleware"
	"go_dfs/internal/model"
	"go_dfs/internal/service"
	"io"
	"net/http"
	"strconv"
	"strings"
	"github.com/google/uuid"
)

type FileHandler struct{
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler{
	return &FileHandler{
		fileService:fileService,
	}
}

const maxUploadSize=100<<20

func validateFilename(name string) error{
	if name==""{
		return errors.New("filename cannot be empty")
	}

	if strings.ContainsAny(name, "\r\n") {
		return errors.New("invalid filename")
	}

	if len(name) > 255 {
		return errors.New("filename too long")
	}

	return nil
}

func safeDownloadFilename(name *string) string{
	if name==nil{
		return "download"
	}
	safeName:=*name
	safeName=strings.ReplaceAll(safeName,"\r","")
	safeName=strings.ReplaceAll(safeName,"\n","")
	safeName=strings.ReplaceAll(safeName,`"`,"")
	return safeName
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request){
	
	r.Body=http.MaxBytesReader(
		w,
		r.Body,
		maxUploadSize,
	)

	userID,ok:=middleware.GetUserID(r)
	if !ok{
		http.Error(w,"unauthorized",http.StatusUnauthorized)
		return
	}

	parsedUserID,err:=uuid.Parse(userID)
	if err!=nil{
		http.Error(w,"invalid user id",http.StatusUnauthorized)
		return
	}

	file,header,err:=r.FormFile("file")
	if err!=nil{
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError){
			http.Error(w,"file too large",http.StatusRequestEntityTooLarge,)
			return
		}	
		http.Error(w,"failed to read uploaded file",http.StatusBadRequest)
		return
	}
	defer file.Close()

	name:=header.Filename
	if err:=validateFilename(name);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	size:=header.Size

	buffer:=make([]byte,512)
	n,err:=file.Read(buffer)
	if err !=nil && err!=io.EOF && err !=io.ErrUnexpectedEOF{
		http.Error(w,"failed to read uploaded file",http.StatusBadRequest)
		return
	}
	contentType:=http.DetectContentType(buffer[:n])

	_,err=file.Seek(0,io.SeekStart)
	if err!=nil{
		http.Error(w,"failed to restart uploaded file",http.StatusInternalServerError)
		return
	}

	uploadedFile,err:=h.fileService.UploadFile(
		r.Context(),
		parsedUserID,
		name,
		size,
		contentType,
		file,
	) 
	if err !=nil{
		http.Error(w,"failed to upload file",http.StatusInternalServerError)
		return
	}

	response:=model.FileResponse{
		ID: uploadedFile.ID,
		Name: uploadedFile.Name,
		Size: uploadedFile.Size,
		ContentType: uploadedFile.ContentType,
		CreatedAt: uploadedFile.CreatedAt,
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	
	json.NewEncoder(w).Encode(response)
}

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request){
	userID,ok:=middleware.GetUserID(r)
	if !ok{
		http.Error(w,"unauthorized",http.StatusUnauthorized)
		return
	}

	parsedUserID,err:=uuid.Parse(userID)
	if err!=nil{
		http.Error(w,"invalid user id",http.StatusUnauthorized)
		return
	}

	fileIDString:=r.PathValue("id")

	fileID,err:=uuid.Parse(fileIDString)
	if err!=nil{
		http.Error(w,"invalid file id",http.StatusBadRequest)
		return
	}

	file,reader,err:=h.fileService.GetFile(
		r.Context(),
		fileID,
		parsedUserID,
	)
	if err!=nil{
		if errors.Is(err,service.ErrFileNotFound){
			http.Error(w,"file not found",http.StatusNotFound)
			return
		}
		http.Error(w,"failed to get file",http.StatusInternalServerError)
		return
	}
	defer reader.Close()
	safeName:=safeDownloadFilename(&file.Name)
	w.Header().Set("Content-Type",file.ContentType)
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment;filename="%s"`,safeName),
	)
	w.Header().Set(
		"Content-Length",
		strconv.FormatInt(file.Size,10),
	)

	_,err=io.Copy(w,reader)
	if err !=nil{
		return
	}
}

func(h *FileHandler) List(w http.ResponseWriter, r *http.Request) {
	userID,ok:=middleware.GetUserID(r)
	if !ok{
		http.Error(w,"unauthorized",http.StatusUnauthorized)
		return
	}

	parsedUserID,err:=uuid.Parse(userID)
	if err!=nil{
		http.Error(w,"invalid user id",http.StatusUnauthorized)
		return
	}

	files,err:=h.fileService.ListFiles(
		r.Context(),
		parsedUserID,
	)
	if err!=nil{
		http.Error(w,"failed to list files",http.StatusInternalServerError)
		return
	}

	responses:=make([]model.FileResponse,0,len(files))
	for _,file:=range files{
		responses=append(responses,model.FileResponse{
			ID: file.ID,
			Name: file.Name,
			Size:file.Size,
			ContentType: file.ContentType,
			CreatedAt: file.CreatedAt,
		})
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responses)
}

func(h *FileHandler) Delete(w http.ResponseWriter, r *http.Request){
	userID,ok:=middleware.GetUserID(r)
	if !ok{
		http.Error(w,"unauthorized",http.StatusUnauthorized)
		return
	}

	parsedUserID,err:= uuid.Parse(userID)
	if err!=nil{
		http.Error(w,"invalid user id",http.StatusUnauthorized)
		return
	}

	filedIDString:=r.PathValue("id")
	fileID,err:=uuid.Parse(filedIDString)
	if err!=nil{
		http.Error(w,"invalid file id",http.StatusBadRequest)
		return
	}

	err=h.fileService.DeleteFile(
		r.Context(),
		fileID,
		parsedUserID,
	)
	if err !=nil{
		if errors.Is(err,service.ErrFileNotFound){
			http.Error(w,"file not found",http.StatusNotFound)
			return
		}
		http.Error(w,"failed to delete file",http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}