package file

import (
	"context"
	"io"
	"log/slog"
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
		slog.Error("Could not load AWS config", "error", err)
		os.Exit(1)
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
		slog.Error("Error parsing s3 url", "error", err)
		os.Exit(1)
	}
	result, err := s.s3Client.GetObject(s.ctx, &s3.GetObjectInput{
		Bucket: aws.String(parsedUrl.Host),
		Key:    aws.String(parsedUrl.Path[1:]),
	})
	if err != nil {
		slog.Error("Error fetching file from s3", "error", err)
		os.Exit(1)
	}
	f, err := os.CreateTemp("/tmp", "target")
	if err != nil {
		slog.Error("failed to create temporary file for s3 target", "error", err)
		os.Exit(1)
	}
	body, err := io.ReadAll(result.Body)
	if err != nil {
		slog.Error("Couldn't read object body.", "error", err)
		os.Exit(1)
	}
	_, err = f.Write(body)
	if err != nil {
		slog.Error("Could not write to temporary file.", "error", err)
		os.Exit(1)
	}
	f.Close()
	return f.Name()
}
