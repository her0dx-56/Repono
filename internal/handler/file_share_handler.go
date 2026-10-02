package handler

import (
	"go_dfs/internal/middleware"
	"go_dfs/internal/service"
	"net/http"
	"go_dfs/internal/model"
	"errors"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type FileShareHandler struct {
	fileShareService *service.FileShareService
}

func NewFileShareHandler(
	fileShareService *service.FileShareService,
) *FileShareHandler{
	return &FileShareHandler{
		fileShareService: fileShareService,
	}
}

func (h *FileShareHandler) Share(w http.ResponseWriter, r *http.Request){
	userID,ok:=middleware.GetUserID(r)
	if !ok{
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ownerID,err:=uuid.Parse(userID)
	if err!=nil{
		http.Error(w,"invalid user id", http.StatusUnauthorized)
		return
	}

	fileID,err:=uuid.Parse(r.PathValue("id"))
	if err!=nil{
		http.Error(w, "invalid file id", http.StatusBadRequest)
		return
	}

	var req model.CreateFileShareRequest

	err=json.NewDecoder(r.Body).Decode(&req)
	if err !=nil{
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.SharedWithUserEmail = strings.TrimSpace(
		req.SharedWithUserEmail,
	)
	if req.SharedWithUserEmail == ""{
		http.Error(
			w,
			"recipient email is required",
			http.StatusBadRequest,
		)
		return
	}

	share,err:=h.fileShareService.CreateShareByEmail(
		r.Context(),
		fileID,
		ownerID,
		req.SharedWithUserEmail,
	)
	if err!=nil{
		switch{
		case errors.Is(err,service.ErrFileNotFound):
			http.Error(w,"file not found", http.StatusNotFound)

		case errors.Is(err,service.ErrUserNotFound):
			http.Error(w,"recipient not found", http.StatusNotFound)

		case errors.Is(err,service.ErrAlreadyShared):
			http.Error(w, "file already shared with user", http.StatusConflict)

		default:
			fmt.Println("CreateShare error:", err)
			http.Error(w, "failed to share file", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(share)
} 

func (h *FileShareHandler) ListSharedFiles(
	w http.ResponseWriter,
	r *http.Request,
){
	userID,ok:=middleware.GetUserID(r)
	if !ok{
		http.Error(w,"unauthorized",http.StatusUnauthorized)
		return
	}

	parsedUserID,err:=uuid.Parse(userID)
	if err !=nil{
		http.Error(w,"invalid user id",http.StatusUnauthorized)
		return
	}

	files,err:=h.fileShareService.GetSharedFiles(
		r.Context(),
		parsedUserID,
	)
	if err!=nil{
		fmt.Println("GetSharedFiles error:", err)
		http.Error(w,"failed to get shared files",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(files)

}