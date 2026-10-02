package handler

import (
	"encoding/json"
	"errors"
	"go_dfs/internal/middleware"
	"go_dfs/internal/service"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type PublicFileShareHandler struct{
	publicFileShareService *service.PublicFileShareService
}

func NewPublicFileShareHandler(
	publicFileShareService *service.PublicFileShareService,
) *PublicFileShareHandler {
	return &PublicFileShareHandler{
		publicFileShareService: publicFileShareService,
	}
}

type CreatePublicFileShareRequest struct{
	ExpiresAt *time.Time `json:"expires_at"`
}

func (h *PublicFileShareHandler) CreatePublicShare(
	w http.ResponseWriter,
	r *http.Request,
){
	userID,ok:=middleware.GetUserID(r)
	if !ok{
		http.Error(w,"unauthorized",http.StatusUnauthorized)
		return
	}
	ownerID,err:=uuid.Parse(userID)
	if err!=nil{
		http.Error(w,"invalid user id",http.StatusUnauthorized)
		return
	}

	fileID,err:=uuid.Parse(r.PathValue("id"))
	if err!=nil{
		http.Error(w, "invalid file id", http.StatusBadRequest)
		return
	}

	var request CreatePublicFileShareRequest

	err=json.NewDecoder(r.Body).Decode(&request)
	if err!=nil{
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	share,err:=h.publicFileShareService.CreatePublicShare(
		r.Context(),
		fileID,
		ownerID,
		request.ExpiresAt,
	)
	if err!=nil{
		if errors.Is(err,service.ErrFileNotFound){
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		http.Error(w,"failed to create public share",http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(share)
}

func(h *PublicFileShareHandler) DownloadPublicFile(
	w http.ResponseWriter,
	r *http.Request,
) {
	token:=r.PathValue("id")
	if token==""{
		http.Error(w, "invalid share token", http.StatusBadRequest)
        return
	}

	file,reader,err:=h.publicFileShareService.GetPublicFile(
		r.Context(),
		token,
	)
	if err!=nil{
		if errors.Is(err,service.ErrFileNotFound){
			http.Error(w, "file not found", http.StatusNotFound)
            return
		}
		http.Error(w,"failed to doenload file",http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type",file.ContentType)
	w.Header().Set("Content-Disposition",`attachment;filename="`+file.Name+`"`)
	w.Header().Set("Content-Length",strconv.FormatInt(file.Size,10),)

	_,err=io.Copy(w,reader)
	if err!=nil{
		return
	}
}	

func (h *PublicFileShareHandler) RevokePublicShare(
	w http.ResponseWriter,
	r *http.Request,
){
	userID,ok:=middleware.GetUserID(r)
	if !ok{
		http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
	}

	ownerID,err:=uuid.Parse(userID)
	if err != nil {
        http.Error(w, "invalid user id", http.StatusUnauthorized)
        return
    }

	fileID,err:=uuid.Parse(r.PathValue("id"))
	 if err != nil {
        http.Error(w, "invalid file id", http.StatusBadRequest)
        return
    }

	err=h.publicFileShareService.RevokePublicShare(
		r.Context(),
		fileID,
		ownerID,
	)
	 if err != nil {
        if errors.Is(err, service.ErrFileNotFound) {
            http.Error(w, "file not found", http.StatusNotFound)
            return
        }

        http.Error(
            w,
            "failed to revoke public share",
            http.StatusInternalServerError,
        )
        return
    }
	w.WriteHeader(http.StatusNoContent)

}