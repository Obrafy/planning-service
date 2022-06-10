package sqs

import (
	"fmt"
	"reflect"

	"github.com/obrafy/planning/infrastructure/constants"
	"github.com/obrafy/planning/infrastructure/sqsbase"
	planningservice "github.com/obrafy/planning/services/planning-service"
	"github.com/sirupsen/logrus"
)

// Messages for /planning path
func PlanningMessageHandler(handlerContext interface{}, msg *sqsbase.SQSMessage, service interface{}) (bool, error) {
	log := logrus.WithFields(logrus.Fields{"incoming-message": msg, "metric": constants.METRIC_PLANNING_HANDLER})

	_, ok := handlerContext.(*MainSQSClient)

	fmt.Printf("Planning message from path %v\n", msg.Path)

	if !ok {
		log.Errorf("can't cast handler context of type %v", reflect.TypeOf(handlerContext))
		return false, fmt.Errorf("can't cast handler context of type %v", reflect.TypeOf(handlerContext))
	}

	planningService, ok := service.(*planningservice.PlanningService)

	if !ok {
		log.Errorf("can't cast service of type %v", reflect.TypeOf(service))
		return false, fmt.Errorf("can't cast service of type %v", reflect.TypeOf(service))
	}

	planningService.Handler(msg) // Call Service Handler

	log.Warn("Incoming Message on PlanningMessageHandler. This is a noop and the message will be deleted without any further action.")

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
