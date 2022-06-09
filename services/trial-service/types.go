package trialservice

import (
	"github.com/obrafy/planning/infrastructure/config"
	"github.com/obrafy/planning/infrastructure/mongobase"
	"github.com/obrafy/planning/infrastructure/s3"
)

type TrialService struct {
	mongobase.MongoServiceBase
	Config          *config.BaseService
	S3ManagerClient *s3.S3ManagerClient
}
