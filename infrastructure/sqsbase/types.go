package sqsbase

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSMessage struct {
	ID            string
	ReceiptHandle string
	Body          string
	Path          string
	Attributes    map[string]string
}

type SQSHandler func(interface{}, *SQSMessage, interface{}) (bool, error)

type HandlerMap map[string]SQSHandler

type ClientSQSBase struct {
	QueueURL           *string
	MessageTimeWaitSec int
	Terminating        bool
	ClientSQS          *sqs.SQS
	HandlerFuncMap     HandlerMap
}

var ErrorNoHandlerProvided = errors.New("no valid handler for received message")
