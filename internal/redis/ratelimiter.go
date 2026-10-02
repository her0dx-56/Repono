package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RateLimiter struct{
	client *Client
}

func NewRateLimiter(client *Client) *RateLimiter{
	return &RateLimiter{
		client:client,
	}
}

func(r *RateLimiter) Allow(
	ctx context.Context,
	key string,
	limit int,
	window time.Duration,
)(bool,error){
	script:=`
		local count = redis.call("INCR",KEYS[1])

		if count==1 then
			redis.call("EXPIRE",KEYS[1],ARGV[1])
		end

		if count<= tonumber(ARGV[2]) then
			return 1
		end

		return 0
	`
	result,err:=r.client.Client.Eval(
		ctx,
		script,
		[]string{key},
		int(window.Seconds()),
		limit,
	).Result()

	if err!=nil{
		return false,err
	}
	return result.(int64)==1,nil
}

func UserRateLimitKey(action string, userID uuid.UUID) string{
	return "rate_limit:" + action + ":user:" + userID.String()
}
