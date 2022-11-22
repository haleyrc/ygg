package r2

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Credentials struct {
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
}

type Client struct {
	s3Client *s3.Client
}

func New(ctx context.Context, creds Credentials) (*Client, error) {
	resolver := makeEndpointResolver(creds.AccountID)
	credsProvider := makeCredentialsProvider(
		creds.AccessKeyID,
		creds.AccessKeySecret,
	)

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(resolver),
		config.WithCredentialsProvider(credsProvider),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := &Client{
		s3Client: s3.NewFromConfig(cfg),
	}

	return client, nil
}

func (c *Client) PresignGetFile(ctx context.Context, bucket, key string) (string, error) {
	client := s3.NewPresignClient(c.s3Client)

	result, err := client.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("presign get file failed: %w", err)
	}

	return result.URL, nil
}

func (c *Client) PutFile(ctx context.Context, bucket, key string, r io.Reader) error {
	fileBytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("put file failed: %w", err)
	}

	contentType := http.DetectContentType(fileBytes)

	_, err = c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("put file failed: %w", err)
	}

	return nil
}

func makeEndpointResolver(accountID string) aws.EndpointResolverWithOptions {
	optionsFunc := func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		endpoint := aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}
		return endpoint, nil
	}
	return aws.EndpointResolverWithOptionsFunc(optionsFunc)
}

func makeCredentialsProvider(accessKeyID, accessKeySecret string) aws.CredentialsProvider {
	return credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")
}
