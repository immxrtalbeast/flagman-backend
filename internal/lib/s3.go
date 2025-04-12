package lib

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadToSupabase(tempPath, fileName, contentType string) (string, error) {
	// Конфигурация S3 клиента
	config := aws.Config{
		Endpoint: aws.String(os.Getenv("SUPABASE_STORAGE_ENDPOINT")),
		Region:   aws.String("eu-north-1"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("SUPABASE_ANON_KEY"),
			os.Getenv("SUPABASE_SECRET_KEY"),
			"",
		),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess, err := session.NewSession(&config)
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(tempPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Загрузка файла
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("SUPABASE_BUCKET_NAME")), // Имя вашего бакета
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	// Формирование публичного URL
	publicURL := fmt.Sprintf("%s/%s/%s",
		os.Getenv("SUPABASE_STORAGE_ENDPOINT"),
		os.Getenv("SUPABASE_BUCKET_NAME"),
		fileName,
	)

	return publicURL, nil
}
