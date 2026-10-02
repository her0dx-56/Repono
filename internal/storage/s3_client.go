package storage

import(
	"context"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(
	ctx context.Context,
	region string,
) (*s3.Client, error){

	cfg,err:=awsconfig. LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(region),
	)
	if err!=nil{
		return nil,err
	}

	return s3.NewFromConfig(cfg),nil
}