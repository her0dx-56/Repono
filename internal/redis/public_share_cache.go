package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"go_dfs/internal/model"

	goredis "github.com/redis/go-redis/v9"
)

type PublicShareCache struct{
	client *Client
}

func NewPublicShareCache (client *Client) *PublicShareCache{
	return &PublicShareCache{
		client:client,
	} 
}

func PublicShareCacheKey(token string) string{
	return "public_share:" + token
}

func (c *PublicShareCache) Get(
	ctx context.Context,
	token string,
) (*model.PublicShare,error) {
	
	key:= PublicShareCacheKey(token)

	data,err:=c.client.Client.Get(ctx,key).Result()
	if err!=nil{
		if err==goredis.Nil{
			return nil,nil
		}
		return nil,err
	}
	
	var share model.PublicShare

	err=json.Unmarshal([]byte(data), &share)
	if err!=nil{
		return nil,err
	}
	return &share,nil
}

func (c *PublicShareCache) Set(
	ctx context.Context,
	share *model.PublicShare,
	ttl time.Duration,
) error{

	if share==nil{
		return errors.New("public share is nil")
	}

	if share.Token == ""{
		return errors.New("public share token is nil")
	}

	data,err:=json.Marshal(share)
	if err !=nil{
		return err
	}

	key:=PublicShareCacheKey(share.Token)

	return c.client.Client.Set(
		ctx,
		key,
		data,
		ttl,
	). Err()
}

func (c *PublicShareCache) Delete(
	ctx context.Context,
	token string,
)error{
	key:=PublicShareCacheKey(token)
	return c.client.Client.Del(ctx,key).Err()
}