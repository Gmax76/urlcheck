package file

import (
	"context"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type (
	S3Controller interface {
		Get(url string) string
	}
	s3Controller struct {
		s3Client *s3.Client
		ctx      context.Context
	}
)

func NewS3Controller(region string) S3Controller {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if region != "" {
		cfg.Region = region
	}
	if err != nil {
		log.Fatalf("Could not load AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	return &s3Controller{
		s3Client: s3Client,
		ctx:      ctx,
	}
}

func (s *s3Controller) Get(filepath string) string {
	parsedUrl, err := url.Parse(filepath)
	if err != nil {
		log.Fatalf("Error parsing s3 url: %v", err)
	}
	result, err := s.s3Client.GetObject(s.ctx, &s3.GetObjectInput{
		Bucket: aws.String(parsedUrl.Host),
		Key:    aws.String(parsedUrl.Path[1:]),
	})
	if err != nil {
		log.Fatalf("Error fetching file from s3: %v", err)
	}
	f, err := os.CreateTemp("/tmp", "target")
	if err != nil {
		log.Fatalf("failed to create temporary file for s3 target, %v", err)
	}
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", parsedUrl.Path[1:], err)
	}
	_, err = f.Write(body)
	if err != nil {
		log.Fatalf("Could not write to temporary file. Error: %v", err)
	}
	f.Close()
	return f.Name()
}
