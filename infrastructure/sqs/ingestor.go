package sqs

import (
	"sync"

	"github.com/obrafy/planning/infrastructure/config"
	"github.com/obrafy/planning/infrastructure/sqsbase"
	"github.com/sirupsen/logrus"
)

type MainSQSClient struct {
	sqsbase.ClientSQSBase
	cfg     *config.SQSClient
	cwMutex *sync.Mutex
}

func NewMainSQSClient(cfg *config.SQSClient) *MainSQSClient {
	return &MainSQSClient{
		sqsbase.ClientSQSBase{
			HandlerFuncMap: HandlerMap,
		},
		cfg,
		&sync.Mutex{},
	}
}

func (msqs *MainSQSClient) Init() error {
	err := msqs.InitSession(
		msqs.cfg.QueueName,
		msqs.cfg.DelaySeconds,
		msqs.cfg.MaximumMessageSize,
		msqs.cfg.MessageRetentionPeriod,
		msqs.cfg.ReceiveMessageWaitTimeSeconds,
		msqs.cfg.AwsRegion,
		msqs.cfg.AwsProfileName,
		msqs.cfg.DlqARN,
		msqs.cfg.MaxReceiveCountDlq,
		msqs.cfg.VisibilityTimeout,
	)

	return err
}

func (msqs *MainSQSClient) Run(services map[string]interface{}) {

	go func() {
		for !msqs.Terminating {
			msqs.cwMutex.Lock()
			err := msqs.ReadAndDispatchMessage(msqs, services)

			if err != nil {
				logrus.Errorf("Error reading SQS: %v", err)
			}

			msqs.cwMutex.Unlock()
		}
	}()
}

func (msqs *MainSQSClient) ShutDown(wg *sync.WaitGroup) {
	msqs.cwMutex.Lock()
	defer msqs.cwMutex.Unlock()
	defer wg.Done()

	msqs.Terminating = true
	logrus.Info("SQS client shutted down")
}
