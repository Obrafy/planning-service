package planningservice

import (
	"github.com/obrafy/planning/infrastructure/config"
	"github.com/obrafy/planning/infrastructure/s3"
)

func NewPlanningService(cfg *config.BaseService, s3ManagerClient *s3.S3ManagerClient) (*PlanningService, error) {
	service := &PlanningService{
		Config:          cfg,
		S3ManagerClient: s3ManagerClient,
	}

	err := service.InitSession(cfg.DatabaseURI, cfg.DatabaseName) // Start Mongo Session for Service

	return service, err
}

func (planningService *PlanningService) Init() error {
	if err := planningService.InitSession(planningService.Config.DatabaseURI, planningService.Config.DatabaseName); err != nil {
		return err
	}

	return nil
}

func (planningService PlanningService) Terminate() {
	planningService.TerminateSession()
}
