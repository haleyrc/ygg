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

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(resolver),
		config.WithCredentialsProvider(makeCredentialsProvider(creds.AccessKeyID, creds.AccessKeySecret)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &Client{s3Client: s3.NewFromConfig(cfg)}, nil
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

	_, err = c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(http.DetectContentType(fileBytes)),
	})
	if err != nil {
		return fmt.Errorf("put file failed: %w", err)
	}
	return nil
}

func makeEndpointResolver(accountID string) aws.EndpointResolverWithOptions {
	optionsFunc := func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		url := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
		return aws.Endpoint{URL: url}, nil
	}
	return aws.EndpointResolverWithOptionsFunc(optionsFunc)
}

func makeCredentialsProvider(accessKeyID, accessKeySecret string) aws.CredentialsProvider {
	return credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")
}

/*
	listObjectsOutput, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range listObjectsOutput.Contents {
		obj, _ := json.MarshalIndent(object, "", "\t")
		fmt.Println(string(obj))
	}
*/
