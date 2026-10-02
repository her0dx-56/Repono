package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(databaseURL string) (*Database ,error) {

	pool,err:= pgxpool.New(context.Background(),databaseURL)
	if(err!=nil){
		return nil, err
	}
	err=pool.Ping(context.Background())
	if err!=nil{
		pool.Close()
		return nil,err
	}

	return &Database{
		Pool:pool,
	},nil

}