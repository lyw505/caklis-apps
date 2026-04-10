package config

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client
var ExternalMinioClient *minio.Client

func ConnectMinio() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	externalEndpoint := os.Getenv("MINIO_EXTERNAL_ENDPOINT")
	if externalEndpoint == "" {
		externalEndpoint = "localhost:9000"
	}
	region := os.Getenv("MINIO_REGION")
	if region == "" {
		region = "us-east-1"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	var err error
	// Internal client for server-to-server communication
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Region: region,
		Secure: useSSL,
	})

	if err != nil {
		log.Fatal("Failed to connect to MinIO:", err)
	}

	// External client for generating presigned URLs for the browser
	ExternalMinioClient, err = minio.New(externalEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Region: region,
		Secure: useSSL,
	})

	if err != nil {
		log.Fatal("Failed to create External MinIO Client:", err)
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

func GetExternalMinioClient() *minio.Client {
	return ExternalMinioClient
}
