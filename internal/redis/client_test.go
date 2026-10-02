package redis


import(
	"context"
	"testing"
)

func TestClient_page(t *testing.T){
	ctx:=context.Background()

	client:=NewClient("localhost:6379")
	defer client.Close()

	err:=client.Ping(ctx)
	if err!=nil{
		t.Fatalf("Ping() error=%v",err)
	}
}