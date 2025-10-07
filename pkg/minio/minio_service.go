package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)



func (m *MinioClient) UploadImage(material string, file io.Reader, fileSize int64, filename string) (string, error) {
    objectName := fmt.Sprintf("%s%s", material, ".png")

    _, err := m.client.PutObject(context.Background(), "materials", objectName, file, fileSize, minio.PutObjectOptions{
        ContentType: "image/png",
    })
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("http://localhost:9000/materials/%s", objectName), nil
}

func (m *MinioClient) DeleteImage(material string) error {
    objectName := fmt.Sprintf("%s.png", material)
    
    err := m.client.RemoveObject(context.Background(), "materials", objectName, minio.RemoveObjectOptions{})
    if err != nil {
        return err
    }
    
    return nil
}
