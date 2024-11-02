package services

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

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
