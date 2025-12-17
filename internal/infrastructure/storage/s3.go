package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewS3Storage(endpoint, region, bucket, accessKey, secretKey, publicURL string) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(_, _ string) (aws.Endpoint, error) {
			return aws.Endpoint{URL: endpoint}, nil
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init S3 config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Storage{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}, nil
}

func (s *S3Storage) PublicURL(key string) string {
	key = strings.TrimLeft(key, "/")
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.publicURL, "/"), key)
}

func (s *S3Storage) Upload(ctx context.Context, key, contentType string, body io.Reader) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}
	return s.PublicURL(key), nil
}

func (s *S3Storage) PresignPut(ctx context.Context, key, contentType string, expiresIn time.Duration) (*PresignedUpload, error) {
	presigner := s3.NewPresignClient(s.client)

	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead, 
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return nil, err
	}

	return &PresignedUpload{
		Key:       key,
		UploadURL: req.URL,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
		ExpiresAt: time.Now().Add(expiresIn),
	}, nil
}

func (s *S3Storage) Copy(ctx context.Context, srcKey, dstKey string) (string, error) {
	copySource := url.PathEscape(s.bucket + "/" + srcKey)

	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(copySource),
		ACL:        types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("copy failed: %w", err)
	}

	return s.PublicURL(dstKey), nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
