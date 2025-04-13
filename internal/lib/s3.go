package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadToSupabase(tempPath, fileName, contentType string) (string, error) {
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

	hash := sha256.Sum256([]byte(fileName))
	hashedString := hex.EncodeToString(hash[:])
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("SUPABASE_BUCKET_NAME")),
		Key:         aws.String(hashedString),
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

func DownloadFromSupabase(fileName, destinationPath string) error {
	// Конфигурация аналогична загрузке
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
		return err
	}

	// Генерация хеша имени файла (как при загрузке)
	hash := sha256.Sum256([]byte(fileName))
	hashedString := hex.EncodeToString(hash[:])

	// Создаем файл для записи
	file, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Используем Downloader для скачивания
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("SUPABASE_BUCKET_NAME")),
			Key:    aws.String(hashedString),
		})

	return err
}
