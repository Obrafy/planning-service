package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/obrafy/planning/infrastructure/config"
)

func NewS3Uploader(config *config.S3UploaderOptions, session *session.Session) *s3manager.Uploader {
	return s3manager.NewUploader(session, func(d *s3manager.Uploader) {
		d.PartSize = *config.DownloaderPartSizeInBytes
		d.Concurrency = *config.Concurrency
	})
}
