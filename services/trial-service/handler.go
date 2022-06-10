package trialservice

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/obrafy/planning/infrastructure/sqsbase"
)

func (trialService *TrialService) Handler(msg *sqsbase.SQSMessage) error {
	// log := logrus.WithFields(logrus.Fields{"incoming-message": msg})

	// Batch Download
	files, err := trialService.S3ManagerClient.Downloader.S3.ListObjects(&s3.ListObjectsInput{
		Bucket: trialService.S3ManagerClient.Bucket,
		Prefix: aws.String("planning-files"),
	})

	if err != nil {
		return fmt.Errorf("error listing files from s3 folder: %w", err)
	}

	for _, file := range files.Contents {
		// Single Object Download

		fileBuffer := new(aws.WriteAtBuffer)

		trialService.S3ManagerClient.Downloader.Download(fileBuffer, &s3.GetObjectInput{
			Bucket: trialService.S3ManagerClient.Bucket,
			Key:    aws.String(*file.Key),
		})

		fmt.Println(string(fileBuffer.Bytes()))

	}

	// fmt.Println(files)

	return nil
}
