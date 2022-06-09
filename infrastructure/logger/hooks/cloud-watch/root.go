package cloudwatchhook

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_BUFFER_SIZE    = 512
	DEFAULT_FLUSH_INTERVAL = 15 * time.Second
	AWS_PROFILE_FILENAME   = ""
)

func (h *CloudWatchHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (h *CloudWatchHook) Fire(entry *logrus.Entry) error {
	message, err := entry.String()

	if err != nil {
		return err
	}

	event := &cloudwatchlogs.InputLogEvent{
		Message:   aws.String(message),
		Timestamp: aws.Int64(time.Now().UnixMilli()),
	}

	go func(ev *cloudwatchlogs.InputLogEvent) {
		h.chLog <- ev
	}(event)

	return err
}

func newAwsConfig(awsRegion, awsProfileName string) *aws.Config {
	return aws.NewConfig().WithRegion(awsRegion).WithCredentials(credentials.NewSharedCredentials(AWS_PROFILE_FILENAME, awsProfileName))
}

func (h *CloudWatchHook) initCloudWatchLogGroup() error {
	resp, err := h.cwl.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(h.groupName),
		LogStreamNamePrefix: aws.String(h.streamName),
	})

	if err != nil {
		if awsError, ok := err.(awserr.Error); ok && awsError.Code() == cloudwatchlogs.ErrCodeResourceNotFoundException {
			if _, err = h.cwl.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
				LogGroupName: aws.String(h.groupName),
			}); err != nil {
				log.Printf("Fail to create group (%s/%s): %v\n", h.groupName, h.streamName, err)
				return err
			}

			return h.initCloudWatchLogGroup()

		} else {
			log.Printf("Fail to get group (%s/%s): %v\n", h.groupName, h.streamName, err)
			return err
		}
	}

	if len(resp.LogStreams) > 0 {
		// We already have a stream / group with this name
		h.nextSequenceToken = resp.LogStreams[0].UploadSequenceToken
	} else {
		// Create a stream if it doesn't exist. the next sequence token will be null
		h.nextSequenceToken = nil
		if _, err := h.cwl.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(h.groupName),
			LogStreamName: aws.String(h.streamName),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (h *CloudWatchHook) writeToCloudWatch(buffer []*cloudwatchlogs.InputLogEvent) {
	if h != nil && len(buffer) > 0 {
		h.cwMutex.Lock()
		defer h.cwMutex.Unlock()

		params := &cloudwatchlogs.PutLogEventsInput{
			LogEvents:     buffer,
			LogGroupName:  aws.String(h.groupName),
			LogStreamName: aws.String(h.streamName),
			SequenceToken: h.nextSequenceToken,
		}

		if resp, err := h.cwl.PutLogEvents(params); err != nil {
			if cloudWatchError, ok := err.(*cloudwatchlogs.InvalidSequenceTokenException); ok && cloudWatchError.ExpectedSequenceToken != nil {
				t := *cloudWatchError.ExpectedSequenceToken
				h.nextSequenceToken = &t
				log.Println("writeToCloudWatch: Error writing to CloudWatch (*cloudwatchlogs.ExpectedSequenceToken). New token:", t)
			} else if _, ok := err.(*cloudwatchlogs.InvalidParameterException); ok {
				s := ""
				for i := 0; i < len(buffer); i++ {
					if buffer[i].Timestamp != nil {
						s += fmt.Sprintf("%d", *buffer[i].Timestamp)
					} else {
						s += "nil"
					}
				}
				log.Println("writeToCloudWatch: Error writing to CloudWatch (*cloudwatchlogs.InvalidParameterException): ", s)
			} else {
				log.Printf("writeToCloudWatch: Error writing to CloudWatch %s (%v)\n", err.Error(), reflect.TypeOf(err))
			}
		} else {
			t := *resp.NextSequenceToken
			h.nextSequenceToken = &t
		}

	}
}

func (h *CloudWatchHook) flushBuffer() {
	if h != nil && len(h.logsBuffer) > 0 {
		h.Lock()
		bufferToFlush := h.logsBuffer
		h.logsBuffer = make([]*cloudwatchlogs.InputLogEvent, 0, h.bufferSize)
		h.Unlock()

		go h.writeToCloudWatch(bufferToFlush)
	}
}

func (h *CloudWatchHook) handleLogs() {
	flushTimer := time.NewTimer(h.flushInterval)
	flushOnExit := make(chan os.Signal)
	signal.Notify(flushOnExit, syscall.SIGINT, syscall.SIGTERM)

	terminate := false

	for !terminate {
		select {
		case le := <-h.chLog:
			h.Lock()
			h.logsBuffer = append(h.logsBuffer, le)
			h.Unlock()

			if len(h.logsBuffer) >= h.flushBufferSize {
				h.flushBuffer()
			}
		case <-flushOnExit:
			h.flushBuffer()
			terminate = true

		case <-flushTimer.C:
			flushTimer.Reset(h.flushInterval)
			h.flushBuffer()
		}

	}
}

func NewCloudWatchHook(groupName, streamName, awsRegion, awsProfileName string, bufferSize int, flushInterval time.Duration) (*CloudWatchHook, error) {

	if bufferSize == 0 {
		bufferSize = DEFAULT_BUFFER_SIZE
	}

	if flushInterval == 0 {
		flushInterval = DEFAULT_FLUSH_INTERVAL
	}

	maxBufferSize := bufferSize + (bufferSize / 4) + 1

	awsConfig, err := session.NewSession(newAwsConfig(awsRegion, awsProfileName))

	if err != nil {
		fmt.Errorf("Cannot initialize Aws Session: %v", err)
	}

	h := &CloudWatchHook{
		cwl:             cloudwatchlogs.New(awsConfig),
		cwMutex:         &sync.Mutex{},
		groupName:       groupName,
		streamName:      streamName,
		chLog:           make(chan *cloudwatchlogs.InputLogEvent),
		logsBuffer:      make([]*cloudwatchlogs.InputLogEvent, 0, maxBufferSize),
		flushBufferSize: bufferSize,
		bufferSize:      maxBufferSize,
		flushInterval:   flushInterval,
	}

	err = h.initCloudWatchLogGroup()

	if err != nil {
		h = nil
	} else {
		go h.handleLogs()
	}

	return h, err
}
