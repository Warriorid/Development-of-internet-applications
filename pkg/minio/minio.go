package minio

import (
	"io"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ImageStorage interface {
    
    UploadImage(material string, file io.Reader, fileSize int64, filename string) (string, error)
    DeleteImage(material string) error
}


type MinioClient struct {
    client *minio.Client
}

func NewMinioClient() (*MinioClient, error) {
    client, err := minio.New("localhost:9000", &minio.Options{
        Creds:  credentials.NewStaticV4("minio", "minio124", ""),
        Secure: false,
    })
    if err != nil {
        return nil, err
    }
    return &MinioClient{client: client}, nil
}