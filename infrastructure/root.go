package infrastructure

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	timehelpers "github.com/obrafy/planning/helpers/time-helpers"
	"github.com/obrafy/planning/infrastructure/config"
	"github.com/obrafy/planning/infrastructure/logger"
	"github.com/obrafy/planning/infrastructure/s3"
	"github.com/obrafy/planning/infrastructure/sqs"
	trialservice "github.com/obrafy/planning/services/trial-service"
	"github.com/sirupsen/logrus"
)

const (
	INIT_METRIC                 = "INIT"
	STARTED_APP_LOG_TIME_FORMAT = "2006/01/02 15:04:05"
)

type App struct {
	Configuration *config.Configuration
	MainSQS       *sqs.MainSQSClient
	MainS3        *s3.S3ManagerClient
	TrialService  *trialservice.TrialService
}

func NewApp(environment string) (*App, error) {
	var (
		app App
		err error
	)

	stopwatch := timehelpers.NewStopWatch()

	cfg, err := config.LoadConfiguration(environment)

	if err != nil {
		log.Println("Error initializing application")
		log.Panicln(err)
	}

	app.Configuration = cfg

	// Logging
	logger.InitLogs(
		app.Configuration.Log.CloudWatchOutput,
		app.Configuration.Log.LogLevel,
		app.Configuration.Log.LogGroup,
		app.Configuration.Log.AwsRegion,
		app.Configuration.Log.AwsProfileName,
		app.Configuration.Log.AwsConfigFile,
		app.Configuration.Log.AwsConfigProfile,
		app.Configuration.Log.AwsStreamPrefix,
		app.Configuration.Log.FileOutput,
		app.Configuration.Log.LogPath,
	)

	// S3 Service
	app.MainS3 = s3.NewS3ManagerClient(&app.Configuration.MainS3Client)

	// Services
	app.TrialService, err = trialservice.NewTrialService(
		&app.Configuration.TrialService,
		app.MainS3,
	)

	if err != nil {
		log.Println("Error initializing application")
		log.Panicln(err)
	}

	if len(app.Configuration.MainSQSQueue.QueueName) > 0 {
		app.MainSQS = sqs.NewMainSQSClient(&app.Configuration.MainSQSQueue)
	}

	logrus.WithField("config", app.Configuration).
		WithField("metric", INIT_METRIC).
		WithField("metric_val", stopwatch.EllapsedMillis()).
		Infof("Started app at %s", time.Now().Format(STARTED_APP_LOG_TIME_FORMAT))

	return &app, err
}

func (app *App) Run() {
	fmt.Println("Starting services")

	if app.MainSQS != nil {
		fmt.Println("Starting SQS")

		if err := app.MainSQS.Init(); err != nil {
			log.Fatalf("Error initializing SQS client %v", err.Error())
		}

		services := make(map[string]interface{})
		services["/trial"] = app.TrialService

		fmt.Println("Running SQS", *app.MainSQS.QueueURL)
		app.MainSQS.Run(services)
	}

	fmt.Println("App Running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Interrupting Application")
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	if app.MainSQS != nil {
		wg.Add(1)
		go app.MainSQS.ShutDown(wg)
	}

	wg.Wait()
}
