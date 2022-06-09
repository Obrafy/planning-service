package trialservice

import (
	"github.com/obrafy/planning/infrastructure/config"
	"github.com/obrafy/planning/infrastructure/s3"
)

func NewTrialService(cfg *config.BaseService, s3ManagerClient *s3.S3ManagerClient) (*TrialService, error) {
	service := &TrialService{
		Config:          cfg,
		S3ManagerClient: s3ManagerClient,
	}

	err := service.InitSession(cfg.DatabaseURI, cfg.DatabaseName) // Start Mongo Session for Service

	return service, err
}

func (trialService *TrialService) Init() error {
	if err := trialService.InitSession(trialService.Config.DatabaseURI, trialService.Config.DatabaseName); err != nil {
		return err
	}

	return nil
}

func (trialService TrialService) Terminate() {
	trialService.TerminateSession()
}
