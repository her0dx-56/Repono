package handler

import (
	"encoding/json"
	"go_dfs/internal/auth"
	"go_dfs/internal/model"
	"go_dfs/internal/service"
	"net/http"
)

type UserHandler struct{
	userService *service.UserService
	jwtManager *auth.JWTManager
}

func NewUserHandler(
	userService *service.UserService,
	jwtManager *auth.JWTManager,
	) *UserHandler{
	return &UserHandler{
		userService:userService,
		jwtManager: jwtManager,
	}
}


func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request){
	var req model.RegisterRequest

	err:= json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		http.Error(w,"invalid request body",http.StatusBadRequest)
		return
	}


	user,err:= h.userService.RegisterUser(r.Context(),req)
	if err !=nil{
		if err == service.ErrEmailAlreadyExists{
			http.Error(w,"email already exist",http.StatusConflict)
			return
		}
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	response := model.UserResponse{
		ID: user.ID.String(),
		Username: user.Username,
		Email: user.Email,
	}
	
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func(h *UserHandler) Login(w http.ResponseWriter, r *http.Request){
	var req model.LoginRequest

	err:= json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		http.Error(w,"invalid credentials",http.StatusBadRequest)
		return
	}
	user,err:=h.userService.LoginUser(r.Context(),req)
	if err !=nil{
		if err == service.ErrInvalidCredentials{
			http.Error(w,"invalid credentials",http.StatusUnauthorized)
			return
		}
		http.Error(w,"failed to login",http.StatusInternalServerError)
	}

	token,err := h.jwtManager.GenerateToken(user.ID.String())
	if err !=nil{
		http.Error(w,"failed to generate token",http.StatusInternalServerError)
		return
	}

	response:=model.LoginResponse{
		Token: token,
		User: model.UserResponse{
			ID: user.ID.String(),
			Username: user.Username,
			Email: user.Email,
		},
	}

	w.Header().Set("Content-Type","appliocation/json")
	json.NewEncoder(w).Encode(response)
}

