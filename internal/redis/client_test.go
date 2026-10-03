package redis

import (
	"context"
	"go_dfs/internal/config"
	"testing"
	"log"
)

func TestClient_page(t *testing.T){
	ctx:=context.Background()

	client,err:=NewClient(config.Load().RedisAddr)
	if err != nil {
   		 log.Fatalf("Redis configuration failed: %v", err)
	}
	defer client.Close()

	err=client.Ping(ctx)
	if err!=nil{
		t.Fatalf("Ping() error=%v",err)
	}
}