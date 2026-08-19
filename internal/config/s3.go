package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	S3Client  *s3.Client
	AWSBucket string
	AWSRegion string
)

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	UsePathStyle    bool
}

func LoadS3() S3Config {
	return S3Config{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:          os.Getenv("AWS_DEFAULT_REGION"),
		Bucket:          os.Getenv("AWS_BUCKET"),
		UsePathStyle:    strings.EqualFold(os.Getenv("AWS_USE_PATH_STYLE_ENDPOINT"), "true"),
	}
}

func InitializeS3() error {
	cfg := LoadS3()
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Bucket == "" {
		return errors.New("s3: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY or AWS_BUCKET env vars are not set")
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return fmt.Errorf("s3: failed to load AWS config: %w", err)
	}
	S3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.UsePathStyle
	})
	AWSBucket = cfg.Bucket
	AWSRegion = cfg.Region
	return nil
}
