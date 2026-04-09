package config

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinio() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatal("Failed to connect to MinIO:", err)
	}

	log.Println("✅ MinIO connected successfully")

	// Ensure bucket exists
	bucketName := os.Getenv("MINIO_BUCKET")
	ctx := context.Background()
	
	exists, err := MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatal("Failed to check bucket:", err)
	}

	if !exists {
		err = MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal("Failed to create bucket:", err)
		}
		log.Printf("✅ Bucket '%s' created successfully\n", bucketName)
	} else {
		log.Printf("✅ Bucket '%s' already exists\n", bucketName)
	}
}

func GetMinioClient() *minio.Client {
	return MinioClient
}
