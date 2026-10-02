package redis


import(
	"context"
	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	Client *goredis.Client
}

func NewClient(addr string) *Client{
	client:=goredis.NewClient(&goredis.Options{
		Addr:addr,
	})

	return &Client{
		Client:client, 
	}

}

func (c *Client) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.Client.Close()
}








