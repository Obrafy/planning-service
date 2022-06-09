package sqsbase

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (sqsb *ClientSQSBase) ReadAndDispatchMessage(handlercontext interface{}, serviceContainer map[string]interface{}) error {
	var ctx context.Context

	timeout := time.Second * time.Duration(sqsb.MessageTimeWaitSec+2)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := sqsb.ClientSQS.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              sqsb.QueueURL,
		MaxNumberOfMessages:   aws.Int64(1),
		WaitTimeSeconds:       aws.Int64(int64(sqsb.MessageTimeWaitSec)),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}

	if len(res.Messages) == 0 {
		return nil
	}

	attrs := make(map[string]string)
	var path string

	for key, attr := range res.Messages[0].MessageAttributes {
		if key == "path" {
			path = *attr.StringValue
		} else {
			attrs[key] = *attr.StringValue
		}
	}

	receivedMessage := &SQSMessage{
		ID:            *res.Messages[0].MessageId,
		ReceiptHandle: *res.Messages[0].ReceiptHandle,
		Body:          *res.Messages[0].Body,
		Attributes:    attrs,
		Path:          path,
	}

	var f SQSHandler = sqsb.HandlerFuncMap["nil"]

	if len(path) > 0 {
		f = sqsb.HandlerFuncMap[path]
	}

	if f == nil {
		f = sqsb.HandlerFuncMap["*"]
	}

	if f != nil {
		delete, err := f(handlercontext, receivedMessage, serviceContainer[path])

		if delete {
			ctx, cancel = context.WithTimeout(ctx, time.Second*2)
			defer cancel()

			_, err := sqsb.ClientSQS.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      sqsb.QueueURL,
				ReceiptHandle: aws.String(receivedMessage.ReceiptHandle),
			})

			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

	} else {
		return ErrorNoHandlerProvided
	}

	return err
}
