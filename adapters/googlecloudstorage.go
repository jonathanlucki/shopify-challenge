package adapters

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

type Storage struct {
	Client *storage.Client
	Bucket *storage.BucketHandle
}

// function initializes storage client
func InitStorage() (*Storage, error) {
	// get bucket name
	bucketName := os.Getenv("STORAGE_BUCKET")
	if bucketName == "" {
		return nil, errors.New("$STORAGE_BUCKET is not set")
	}

	// establish client connection
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	// connect to storage bucket
	bucket := client.Bucket(bucketName)

	storage := Storage{client, bucket}
	return &storage, nil
}

// function returns url of image with $id in storage bucket
// it does not check that the image with $id exists in storage
func GetImageUrl(id string) (string, error) {
	// get base url
	baseUrl := os.Getenv("STORAGE_URL")
	if baseUrl == "" {
		return "", errors.New("$STORAGE_URL is not set")
	}

	// construct and return url
	url := baseUrl + id + ".png"
	return url, nil
}

// function uploads image file to storage bucket
func (s *Storage) UploadImage(c *gin.Context, id string, file *multipart.FileHeader) error {
	fileName := id + ".png"

	// open file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// open bucket
	dst := s.Bucket.Object(fileName).NewWriter(c)
	defer dst.Close()

	// copy file into bucket
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

// function deletes image file to storage bucket
func (s *Storage) DeleteImage(c *gin.Context, id string) error {
	fileName := id + ".png"

	if err := s.Bucket.Object(fileName).Delete(c); err != nil {
		return err
	}

	return nil
}