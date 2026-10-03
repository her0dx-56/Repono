package server

import (
	"fmt"
	"go_dfs/internal/auth"
	"go_dfs/internal/config"
	"go_dfs/internal/database"
	"go_dfs/internal/handler"
	"go_dfs/internal/middleware"
	"go_dfs/internal/redis"
	"go_dfs/internal/repository"
	"go_dfs/internal/service"
	"go_dfs/internal/storage"
	"net/http"
	"time"
	"context"
	"strings"
	
)
type Server struct{
	Config config.Config
	Router *http.ServeMux
	DB   *database.Database
	RateLimiter *redis.RateLimiter
}

func New(
	cfg config.Config, 
	db *database.Database,
	redisClient *redis.Client,
	) (*Server,error){
	router:=http.NewServeMux()

	rateLimiter:=redis.NewRateLimiter(redisClient)

	jwtManager:=auth.NewJWTManager(cfg.JWTSecret)
	authTestHandler:=middleware.Auth(jwtManager)(
		http.HandlerFunc(handler.AuthTest),
	)
	router.Handle(
		"GET /auth-test",
		authTestHandler,
	)
	
	s3Client,err:=storage.NewS3Client(
		context.Background(),
		cfg.AWSRegion,
	)
	if err!=nil{
		return nil,fmt.Errorf("failed to create s3 client:%w",err)
	}
	s3Storage:=storage.NewS3Storage(
		s3Client,
		cfg.AWSS3Bucket,
	)
	
	fileShareRepository:=repository.NewFileShareRepository(db.Pool)

	userRepository:=repository.NewUserRepository(db.Pool)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService,jwtManager)

	healthService:=service.NewHealthService()
	healthHandler:=handler.NewHealthHandler(healthService)

	fileRepository:=repository.NewFileRepository(db.Pool)
	// localStorage:=storage.NewLocalStorage("./storage")
	fileService:=service.NewFileService(
		fileRepository,
		fileShareRepository,
		// localStorage,
		s3Storage,
	)
	fileHandler:=handler.NewFileHandler(fileService)

	router.Handle(
		"POST /files",
		middleware.Auth(jwtManager)(
			middleware.RateLimit(
				rateLimiter,
				"upload",
				5,
				time.Minute,
			)(
				http.HandlerFunc(fileHandler.Upload),
			),
		),
	)

	router.Handle(
		"GET /files",
		middleware.Auth(jwtManager)(
			http.HandlerFunc(fileHandler.List),
		),
	)


	router.Handle(
		"GET /files/{id}",
		middleware.Auth(jwtManager)(
			http.HandlerFunc(fileHandler.Download),
		),
	)

	router.Handle(
		"DELETE /files/{id}",
		middleware.Auth(jwtManager)(
			http.HandlerFunc(fileHandler.Delete),
		),
	)

	
	fileShareService:=service.NewFileShareService(
		fileRepository,
		fileShareRepository,
		userRepository,
	)
	fileShareHandler:=handler.NewFileShareHandler(fileShareService)
	router.Handle(
		"POST /files/{id}/shares",
		middleware.Auth(jwtManager)(
			http.HandlerFunc(fileShareHandler.Share),
		),
	)
	
	router.Handle(
		"GET /files/shared",
		middleware.Auth(jwtManager)(
			http.HandlerFunc(fileShareHandler.ListSharedFiles),
		),
	)	
	publicShareCache:=redis.NewPublicShareCache(redisClient)

	publicFileShareRepository:=repository.NewPublicFileShareRepository(db.Pool)
	publicFileShareService:=service.NewPublicFileShareService(
		fileRepository,
		publicFileShareRepository,
		publicShareCache,
		// localStorage,
		s3Storage,
	)
	publicFileShareHandler:=handler.NewPublicFileShareHandler(
		publicFileShareService,
	)
	router.Handle(
		"POST /files/{id}/public-share",
		middleware.Auth(jwtManager)(
			http.HandlerFunc(publicFileShareHandler.CreatePublicShare),
		),
	)
	router.Handle(
		"GET /share/{id}",
		http.HandlerFunc(publicFileShareHandler.DownloadPublicFile),
	)
	router.Handle(
		"DELETE /files/{id}/public-share",
		middleware.Auth(jwtManager)(
			http.HandlerFunc(publicFileShareHandler.RevokePublicShare),
		),
	)

	router.HandleFunc("GET /health",healthHandler.Health)
	router.HandleFunc("POST /users/register",userHandler.Register)
	router.HandleFunc("POST /users/login",userHandler.Login)

	return &Server{
		Config:cfg,
		Router: router,
		DB:db,
		RateLimiter: rateLimiter,
	},nil
}

func(s*Server) Start() error{
	address:= fmt.Sprintf("%s:%s",s.Config.Host, s.Config.Port)
	fmt.Println("Server starting on:",address)

	corsHandler:=http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		origin:=r.Header.Get("Origin")

		configuredOrigin:=strings.TrimRight(
			strings.TrimSpace(s.Config.FrontendOrigin),"/",
		)
		allowedOrigin := origin == "http://localhost:5173" ||
    		origin == configuredOrigin
		 if allowedOrigin{
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if r.Method==http.MethodOptions{
			w.WriteHeader(http.StatusNoContent)
			return
		}
		s.Router.ServeHTTP(w,r)
	})
	return http.ListenAndServe(address,corsHandler)
}
