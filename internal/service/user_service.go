package service

import( 
	"go_dfs/internal/repository"
	"context"
	"errors"
	"go_dfs/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyExists =errors.New("email already exists")
var ErrInvalidCredentials=errors.New("invalid credentials")


type UserService struct{
	userRepository *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

func(s *UserService) RegisterUser(
	ctx context.Context, 
	req model.RegisterRequest,
) (*model.User,error) {

	existingUser,err:=s.userRepository.GetUserByEmail(ctx,req.Email )
	if err !=nil{
		return nil,err
	}

	if existingUser !=nil{
		return nil, ErrEmailAlreadyExists
	}
	
	hashedPassword,err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err !=nil{
		return nil, err
	}
	user := &model.User{
		Username: req.Username,
		Email: req.Email,
		PasswordHash: string(hashedPassword),
	}
	err=s.userRepository.CreateUser(ctx,user)
	if err !=nil{
		return nil,err
	}
	return user,nil
}

func (s *UserService) LoginUser(
	ctx  context.Context,
	req model.LoginRequest,
) (*model.User,error){
	user,err:= s.userRepository.GetUserByEmail(ctx,req.Email)
	if err != nil{
		return nil,err
	}
	if user==nil{
		return nil,ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err!= nil{
		return nil,ErrInvalidCredentials
	}
	return user,nil
}