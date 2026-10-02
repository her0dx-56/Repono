package storage

import(
	"context"
	"io"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/aws"

)

type S3Storage struct{
	client *s3.Client
	bucket string
}

func NewS3Storage(
	client *s3.Client,
	bucket string,
) *S3Storage {
	return &S3Storage{
		client:client,
		bucket:bucket,
	}	
}

func (s *S3Storage) Save(
	ctx context.Context,
	key string,
	r io.Reader,
) error{
	_,err:=s.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:aws.String(key),
			Body:r,
		},
	)
	return err
}

func(s *S3Storage) Get(
	ctx context.Context,
	key string,
) (io.ReadCloser,error) {
	result,err:=s.client.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket:aws.String(s.bucket),
			Key:aws.String(key),
		},
	)
	if err!=nil{
		return nil,err
	}
	return result.Body,nil
}

func(s *S3Storage) Delete(
	ctx context.Context,
	key string,
) error {
	_,err:=s.client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket:aws.String(s.bucket),
			Key:aws.String(key),
		},
	)
	return err
}
