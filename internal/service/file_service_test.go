package service

import (
	"context"
	"go_dfs/internal/repository"
	"go_dfs/internal/storage"
	"os"
	"testing"
	"strings"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestFileService_UploadFile(t *testing.T){
	ctx:=context.Background()
	databaseURL:=os.Getenv("TEST_DATABASE_URL")
	if databaseURL==""{
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	pool,err:=pgxpool.New(ctx,databaseURL)
	if err!=nil{
		t.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	err=pool.Ping(ctx)
	if err!=nil{
		t.Fatalf("failed to connect to test database:%v",err)
	}

	fileRepository:=repository.NewFileRepository(pool)
	fileShareRepository:=repository.NewFileShareRepository(pool)
	dir:=t.TempDir()
	localStorage:=storage.NewLocalStorage(dir)

	fileService:=NewFileService(
		fileRepository,
		fileShareRepository,
		localStorage,
	)

	userID:=uuid.New()

	_,err=pool.Exec(
		ctx,
		`
		INSERT INTO users(
			id,
			username,
			email,
			password_hash
		)
		VALUES($1,$2,$3,$4)
		`,
		userID,
		"test-user",
		userID.String()+"@example.com",
		"test-password-hash",
	)
	if err!=nil{
		t.Fatalf("failed to create test user:%v",err)
	}

	content:="hello from service test"
	reader:=strings.NewReader(content)

	_,err=fileService.UploadFile(
		ctx,
		userID,
		"test.txt",
		int64(len(content)),
		"text/plain",
		reader,
	)
	if err!=nil{
		t.Fatalf("UploadFile() error = %v", err)
	}

} 