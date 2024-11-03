package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/davidalvarez305/yd_vending/constants"
)

func UploadFileToS3(file multipart.File, fileSize int64, s3FilePath string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.AWSRegion),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	buffer := make([]byte, fileSize)
	file.Read(buffer)
	fileType := http.DetectContentType(buffer)

	svc := s3.New(sess)

	input := &s3.PutObjectInput{
		Bucket:        aws.String(constants.AWSS3BucketName),
		Key:           aws.String(s3FilePath),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String(fileType),
	}

	_, err = svc.PutObject(input)
	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %w", err)
	}

	return nil
}

func DownloadFileFromS3(s3FilePath, localFilePath string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(constants.AWSRegion),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(constants.AWSS3BucketName),
		Key:    aws.String(s3FilePath),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return "", fmt.Errorf("failed to download file from S3: %w", err)
	}
	defer result.Body.Close()

	outFile, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer outFile.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(io.MultiWriter(outFile, &buffer), result.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy file to local file: %w", err)
	}

	fmt.Printf("File downloaded successfully to %s\n", localFilePath)
	return localFilePath, nil
}
