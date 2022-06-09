package sqsbase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	DefaultMessageTimeWaitSec = 30
	DefaultAwsRegion          = endpoints.SaEast1RegionID
	DefaultAwsProfileName     = "default"
	DefaultMaxReceiveCount    = "3"
	AWS_PROFILE_FILENAME      = ""
)

func (sqsb *ClientSQSBase) InitSession(
	queueName string,
	delaySeconds *int,
	maximumMessageSize *int,
	messageRetentionPeriod *int,
	receiveMessageWaitTimeSeconds *int,
	awsRegion *string,
	awsProfileName *string,
	dlqARN *string,
	maxReceiveCountDlq *int,
	visibilityTimeout *int,
) (err error) {

	if receiveMessageWaitTimeSeconds != nil {
		sqsb.MessageTimeWaitSec = *receiveMessageWaitTimeSeconds
	} else {
		sqsb.MessageTimeWaitSec = DefaultMessageTimeWaitSec
	}

	var region, profileName string

	if awsRegion != nil {
		region = *awsRegion
	} else {
		region = DefaultAwsRegion
	}

	if awsProfileName != nil {
		profileName = *awsProfileName
	} else {
		profileName = DefaultAwsProfileName
	}

	session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials(AWS_PROFILE_FILENAME, profileName),
	}))

	sqsb.ClientSQS = sqs.New(session)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attributes := make(map[string]*string)

	if delaySeconds != nil {
		attributes["DelaySeconds"] = aws.String(strconv.Itoa(*delaySeconds))
	}

	if maximumMessageSize != nil {
		attributes["MaximumMessageSize"] = aws.String(strconv.Itoa(*maximumMessageSize))
	}

	if messageRetentionPeriod != nil {
		attributes["MessageRetentionPeriod"] = aws.String(strconv.Itoa(*messageRetentionPeriod))
	}

	attributes["ReceiveMessageWaitTimeSeconds"] = aws.String(strconv.Itoa(sqsb.MessageTimeWaitSec))

	if visibilityTimeout != nil {
		attributes["VisibilityTimeout"] = aws.String(strconv.Itoa(*visibilityTimeout))
	}

	if dlqARN != nil {
		var maxReceiveCount string

		if maxReceiveCountDlq != nil {
			maxReceiveCount = strconv.Itoa(*maxReceiveCountDlq)
		} else {
			maxReceiveCount = DefaultMaxReceiveCount
		}

		dlqPolicy, err := json.Marshal(map[string]string{
			"deadLetterTargetArn": *dlqARN,
			"maxReceiveCount":     maxReceiveCount,
		})

		if err != nil {
			return fmt.Errorf("error generating DLQ attributes")
		}

		attributes[sqs.QueueAttributeNameRedrivePolicy] = aws.String(string(dlqPolicy))

	}

	if len(attributes) == 0 {
		attributes = nil
	}

	queueParams := sqs.CreateQueueInput{
		QueueName:  &queueName,
		Attributes: attributes,
	}

	out, err := sqsb.ClientSQS.CreateQueueWithContext(ctx, &queueParams)

	if err == nil {
		sqsb.QueueURL = out.QueueUrl
		return nil
	}

	return err

}
