package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/obrafy/planning/infrastructure/config"
)

const (
	DefaultAwsRegion      = endpoints.SaEast1RegionID
	DefaultAwsProfileName = "default"
	AWS_PROFILE_FILENAME  = ""
)

func NewS3ManagerClient(config *config.S3Client) *S3ManagerClient {
	var region, profileName string

	if config.AwsRegion != nil {
		region = *config.AwsRegion
	} else {
		region = DefaultAwsRegion
	}

	if config.AwsProfileName != nil {
		profileName = *config.AwsProfileName
	} else {
		profileName = DefaultAwsProfileName
	}

	// New Session
	session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials(AWS_PROFILE_FILENAME, profileName),
	}))

	dowloader := NewS3Downloader(config.DownloaderOptions, session)
	uploader := NewS3Uploader(config.UploaderOptions, session)

	return &S3ManagerClient{
		Bucket:     config.BucketName,
		Downloader: dowloader,
		Uploader:   uploader,
	}
}
