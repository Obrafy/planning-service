package sqs

import (
	"github.com/obrafy/planning/infrastructure/constants"
	"github.com/obrafy/planning/infrastructure/sqsbase"
)

var HandlerMap = sqsbase.HandlerMap{
	// ... Other Handlers
	constants.INCOMING_MESSAGE_PATH_PLANNING:  PlanningMessageHandler,
	constants.INCOMING_MESSAGE_PATH_CATCH_ALL: GenericHandler, // Generic Handler for when no path is a match
}
