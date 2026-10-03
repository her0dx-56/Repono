package redis


import(
	"context"
	"fmt"
	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	Client *goredis.Client
}

func NewClient(redisURL string) (*Client,error){
	options,err:=goredis.ParseURL(redisURL)
	if err !=nil{
		return nil, fmt.Errorf("invalid Redis URL:%w",err)
	}
	client:=goredis.NewClient(options)
	

	return &Client{
		Client:client, 
	},nil

}

func (c *Client) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.Client.Close()
}








