package helpers

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

const maxImageSizeBytes = 5 << 20

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

type UploadResult struct {
	Key      string
	URL      string
	FileName string
}

func validateImage(fh *multipart.FileHeader) error {
	if fh.Size > maxImageSizeBytes {
		return fmt.Errorf("%s: %w", fh.Filename, constants.ImageTooLarge)
	}
	contentType := fh.Header.Get("Content-Type")
	if !allowedImageTypes[strings.ToLower(contentType)] {
		return fmt.Errorf("%s: %w", fh.Filename, constants.InvalidImageType)
	}
	return nil
}

func buildKey(folder, fileName string) string {
	ext := filepath.Ext(fileName)
	folder = strings.Trim(folder, "/")
	if folder == "" {
		folder = "uploads"
	}
	return fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)
}

func UploadSingleImage(ctx context.Context, fh *multipart.FileHeader, folder string) (*UploadResult, error) {
	if config.S3Client == nil {
		return nil, errors.New("s3: client not initialized, call config.InitializeS3() at startup")
	}
	if err := validateImage(fh); err != nil {
		return nil, err
	}

	file, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("s3: failed to open uploaded file %s: %w", fh.Filename, err)
	}
	defer file.Close()

	key := buildKey(folder, fh.Filename)
	contentType := fh.Header.Get("Content-Type")

	_, err = config.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(config.AWSBucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", constants.ImageUploadFailed, err.Error())
	}

	return &UploadResult{Key: key, URL: publicURL(key), FileName: fh.Filename}, nil
}

func UploadMultipleImage(ctx context.Context, files []*multipart.FileHeader, folder string) ([]UploadResult, error) {
	if len(files) == 0 {
		return nil, errors.New("s3: no files provided")
	}

	results := make([]UploadResult, 0, len(files))
	var errs []error

	for _, fh := range files {
		res, err := UploadSingleImage(ctx, fh, folder)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results = append(results, *res)
	}

	if len(errs) > 0 {
		return results, errors.Join(errs...)
	}
	return results, nil
}

func DeleteImage(ctx context.Context, key string) error {
	if config.S3Client == nil {
		return errors.New("s3: client not initialized, call config.InitializeS3() at startup")
	}
	if key == "" {
		return errors.New("s3: key is required")
	}

	_, err := config.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(config.AWSBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("%w: %s", constants.ImageDeleteFailed, err.Error())
	}
	return nil
}

func GetImage(ctx context.Context, key string, expiry time.Duration) (string, error) {
	if config.S3Client == nil {
		return "", errors.New("s3: client not initialized, call config.InitializeS3() at startup")
	}
	if key == "" {
		return "", errors.New("s3: key is required")
	}
	if expiry <= 0 {
		expiry = 15 * time.Minute
	}

	presignClient := s3.NewPresignClient(config.S3Client)
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(config.AWSBucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("%w: %s", constants.ImageNotFound, err.Error())
	}
	return req.URL, nil
}

func publicURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.AWSBucket, config.AWSRegion, key)
}
