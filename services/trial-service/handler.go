package trialservice

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	batchDownloadObject := []s3manager.BatchDownloadObject{}
	buffers := []*aws.WriteAtBuffer{}

	for _, file := range files.Contents {
		fileBuffer := new(aws.WriteAtBuffer)

		batchDownloadObject = append(
			batchDownloadObject,
			s3manager.BatchDownloadObject{
				Object: &s3.GetObjectInput{
					Bucket: trialService.S3ManagerClient.Bucket,
					Key:    aws.String(*file.Key),
				},
				Writer: fileBuffer,
			},
		)

		buffers = append(buffers, fileBuffer)
	}

	batchDownloadIterator := &s3manager.DownloadObjectsIterator{Objects: batchDownloadObject}

	if err := trialService.S3ManagerClient.Downloader.DownloadWithIterator(
		context.Background(),
		batchDownloadIterator,
	); err != nil {
		return fmt.Errorf("error downloading files from s3 folder: %w", err)
	}

	for _, buffer := range buffers {
		fmt.Println(string(buffer.Bytes()))
	}

	// fmt.Println(files)

	return nil
}
