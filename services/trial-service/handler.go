package trialservice

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/obrafy/planning/infrastructure/sqsbase"
	"github.com/sirupsen/logrus"
)

func (trialService *TrialService) Handler(msg *sqsbase.SQSMessage) error {
	log := logrus.WithFields(logrus.Fields{"incoming-message": msg})

	fileBuffer := new(aws.WriteAtBuffer)

	trialService.S3ManagerClient.Downloader.Download(fileBuffer, &s3.GetObjectInput{
		Bucket: trialService.S3ManagerClient.Bucket,
		Key:    aws.String("planning-files/new_user_credentials-3.csv"),
	})

	fmt.Println(string(fileBuffer.Bytes()))

	log.Info("New Trial Message")
	fmt.Println("New Trial Message In Service Handler")

	return nil
}
