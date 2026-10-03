package main

import (
	"context"
	"go_dfs/internal/config"
	"go_dfs/internal/database"
	"go_dfs/internal/redis"
	"go_dfs/internal/server"
	"log"
)

func main(){
	cfg:=config.Load()
	db,err:=database.New(cfg.DatabaseURL)
	if err!=nil{
		log.Fatal("Database connection failed:",err)
	}

	redisClient,err:=redis.NewClient(cfg.RedisAddr)
	if err != nil {
    	log.Fatalf("Redis configuration failed: %v", err)
	}
	defer redisClient.Close()

	err=redisClient.Ping(context.Background())
	if err!=nil{
		log.Fatal("Redis connection failed:",err)
	}

	srv,err:=server.New(cfg,db,redisClient)
	if err!=nil{
		log.Fatal("server initialization failed:",err)
	}

	err=srv.Start()
	if err!=nil{
		log.Fatal(err)
	}



	
}