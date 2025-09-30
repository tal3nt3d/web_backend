package minio

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"web_backend/internal/app/ds"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(endpoint, accessKey, secretKey string, useSSL bool) (*minio.Client, error) {
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
}

func UploadImage(ctx context.Context, client *minio.Client, bucket string, file *multipart.FileHeader, device ds.Device) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	contentType := file.Header.Get("Content-Type")
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		default:
			ext = ".bin"
		}
	}
	objectName := fmt.Sprintf("%s%s", device.Title, ext)

	_, err = UploadFromReader(ctx, client, bucket, objectName, f, file.Size, contentType)
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func UploadFromReader(ctx context.Context, client *minio.Client, bucket, objectName string, r io.Reader, size int64, contentType string) (minio.UploadInfo, error) {
	info, err := client.PutObject(ctx, bucket, objectName, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return info, err
}

func DeleteObject(ctx context.Context, client *minio.Client, bucket, objectName string) error {
	return client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

func InitMinio() (*minio.Client, error) {
	host := os.Getenv("MINIO_HOST")
	port := os.Getenv("MINIO_PORT")
	user := os.Getenv("MINIO_USER")
	pass := os.Getenv("MINIO_PASS")
	mc, err := NewMinioClient(host+":"+port, user, pass, false)
	if err != nil {
		return nil, err
	}
	return mc, nil
}

func GetImgBucket() string {
	return "test"
}
