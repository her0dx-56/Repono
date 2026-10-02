package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct{
	Host string
	Port string
	DatabaseURL string
	JWTSecret string
	RedisAddr   string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSS3Bucket        string
}

func Load()Config{
	_=godotenv.Load()
	return Config{
		Host:"localhost",
		Port:"9000",
		DatabaseURL:os.Getenv("DATABASE_URL") ,
		JWTSecret: os.Getenv("JWT_SECRET"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          os.Getenv("AWS_REGION"),
		AWSS3Bucket:        os.Getenv("AWS_S3_BUCKET"),
	}
}