package sqs

import (
	"fmt"
	"reflect"

	"github.com/obrafy/planning/infrastructure/constants"
	"github.com/obrafy/planning/infrastructure/sqsbase"
	trialservice "github.com/obrafy/planning/services/trial-service"
	"github.com/sirupsen/logrus"
)

// Messages for /trial path
func TrialMessageHandler(handlerContext interface{}, msg *sqsbase.SQSMessage, service interface{}) (bool, error) {
	log := logrus.WithFields(logrus.Fields{"incoming-message": msg, "metric": constants.METRIC_TRIAL_HANDLER})

	_, ok := handlerContext.(*MainSQSClient)

	fmt.Printf("Trial message from path %v\n", msg.Path)

	if !ok {
		log.Errorf("can't cast handler context of type %v", reflect.TypeOf(handlerContext))
		return false, fmt.Errorf("can't cast handler context of type %v", reflect.TypeOf(handlerContext))
	}

	trialService, ok := service.(*trialservice.TrialService)

	if !ok {
		log.Errorf("can't cast service of type %v", reflect.TypeOf(service))
		return false, fmt.Errorf("can't cast service of type %v", reflect.TypeOf(service))
	}

	trialService.Handler(msg) // Call Service Handler

	log.Warn("Incoming Message on TrialMessageHandler. This is a noop and the message will be deleted without any further action.")

	return true, nil
}

// Handle Messages with no match for path
func GenericHandler(handlerContext interface{}, msg *sqsbase.SQSMessage, service interface{}) (bool, error) {
	log := logrus.WithFields(logrus.Fields{"incoming-message": msg, "metric": constants.METRIC_GENERIC_HANDLER})

	_, ok := handlerContext.(*MainSQSClient)

	fmt.Printf("Generic message from path %v\n", msg.Path)

	if !ok {
		log.Errorf("can't cast handler context of type %v", reflect.TypeOf(handlerContext))
		return false, fmt.Errorf("can't cast handler context of type %v", reflect.TypeOf(handlerContext))
	}

	log.Warn("Incoming Message on GenericRouteHandler. This is a noop and the message will be deleted without any further action.")

	return true, nil
}
