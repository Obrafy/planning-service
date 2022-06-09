package config

type Configuration struct {
	MainSQSQueue SQSClient   `yaml:"main-sqs-client"`
	MainS3Client S3Client    `yaml:"main-s3-client"`
	TrialService BaseService `yaml:"trial-service"`
	ObrafyAPI    APIClient   `yaml:"obrafy-api"`
	Log          Log         `yaml:"log"`
}

type SQSClient struct {
	QueueName                     string  `yaml:"queue-name"`
	DelaySeconds                  *int    `yaml:"delay-seconds"`
	MaximumMessageSize            *int    `yaml:"maximum-message-size"`
	MessageRetentionPeriod        *int    `yaml:"message-retention-period"`
	ReceiveMessageWaitTimeSeconds *int    `yaml:"receive-message-wait-time-seconds"`
	AwsRegion                     *string `yaml:"aws-region"`
	AwsProfileName                *string `yaml:"aws-profile-name"`
	DlqARN                        *string `yaml:"dlq-arn"`
	MaxReceiveCountDlq            *int    `yaml:"max-receive-count-dlq"`
	VisibilityTimeout             *int    `yaml:"visibility-timeout"`
}

type S3DownloaderOptions struct {
	DownloaderPartSizeInBytes *int64 `yaml:"dowloader-part-size-in-bytes"`
	Concurrency               *int   `yaml:"concurrency"`
}

type S3UploaderOptions struct {
	DownloaderPartSizeInBytes *int64 `yaml:"dowloader-part-size-in-bytes"`
	Concurrency               *int   `yaml:"concurrency"`
}

type S3Client struct {
	BucketName        *string              `yaml:"aws-bucket-name"`
	AwsRegion         *string              `yaml:"aws-region"`
	AwsProfileName    *string              `yaml:"aws-profile-name"`
	DownloaderOptions *S3DownloaderOptions `yaml:"downloader-options"`
	UploaderOptions   *S3UploaderOptions   `yaml:"uploader-options"`
}

type BaseService struct {
	DatabaseURI  string `yaml:"database-uri"`
	DatabaseName string `yaml:"database-name"`
	Collection   string `yaml:"collection"`
}

type APIClient struct {
	BaseURI          string `yaml:"base-uri"`
	TimeoutInSeconds int    `yaml:"timeout-in-seconds"`
}

type Log struct {
	CloudWatchOutput bool   `yaml:"aws-cloud-watch-active"`
	LogLevel         string `yaml:"log-level"`
	LogGroup         string `yaml:"log-group"`
	AwsRegion        string `yaml:"aws-region"`
	AwsProfileName   string `yaml:"aws-profile-name"`
	AwsConfigFile    string `yaml:"aws-config-file"`
	AwsConfigProfile string `yaml:"aws-config-profile"`
	AwsStreamPrefix  string `yaml:"aws-stream-prefix"`
	FileOutput       bool   `yaml:"file-active"`
	LogPath          string `yaml:"log-path"`
}
