package redis

import (
	"go_dfs/internal/config"
	"testing"
	"log"
	"context"
	"time"

	"github.com/google/uuid"
)

func TestRateLimiter_Allow(t *testing.T){
	ctx:=context.Background()

	client,err:=NewClient(config.Load().RedisAddr)
	if err!=nil{
		log.Fatalf("Redis configuration failed: %v", err)
	}
	defer client.Close()
	
	rateLimiter:=NewRateLimiter(client)

	key:="test:rate-limit:" + uuid.New().String()

	t.Cleanup(func(){
	client.Client.Del(ctx,key)
	})

	for i := 1; i <= 4; i++ {
		allowed, err := rateLimiter.Allow(
			ctx,
			key,
			3,
			10*time.Second,
		)
		ttl, err := client.Client.TTL(ctx, key).Result()

		if err != nil {
			t.Fatalf("TTL() error = %v", err)
		}

		if ttl <= 0 || ttl > 10*time.Second {
			t.Fatalf("unexpected TTL: %v", ttl)
		}

		if err != nil {
			t.Fatalf("Allow() request %d error = %v", i, err)
		}

		if i <= 3 && !allowed {
			t.Fatalf("request %d should be allowed", i)
		}

		if i == 4 && allowed {
			t.Fatalf("request %d should be rejected", i)
		}
	}
}

