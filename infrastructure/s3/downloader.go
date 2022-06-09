package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/obrafy/planning/infrastructure/config"
)

func NewS3Downloader(config *config.S3DownloaderOptions, session *session.Session) *s3manager.Downloader {
	return s3manager.NewDownloader(session, func(d *s3manager.Downloader) {
		d.PartSize = *config.DownloaderPartSizeInBytes
		d.Concurrency = *config.Concurrency
	})
}
