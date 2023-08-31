package minio

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Must инициализирует minio-клиента.
func Must(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
) *minio.Client {
	client, err := minio.New(
		endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: false,
		},
	)

	if err != nil {
		panic(fmt.Errorf("new minio client: %v", err))
	}

	return client
}
