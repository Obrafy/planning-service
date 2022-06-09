package cloudwatchhook

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type CloudWatchHook struct {
	sync.Mutex
	cwl               *cloudwatchlogs.CloudWatchLogs
	cwMutex           *sync.Mutex
	groupName         string
	streamName        string
	nextSequenceToken *string
	bufferSize        int
	flushBufferSize   int
	flushInterval     time.Duration
	chLog             chan *cloudwatchlogs.InputLogEvent
	logsBuffer        []*cloudwatchlogs.InputLogEvent
}
