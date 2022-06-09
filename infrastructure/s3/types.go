package s3

import "github.com/aws/aws-sdk-go/service/s3/s3manager"

type S3ManagerClient struct {
	Bucket     *string
	Downloader *s3manager.Downloader
	Uploader   *s3manager.Uploader
}
