package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	cloudwatchhook "github.com/obrafy/planning/infrastructure/logger/hooks/cloud-watch"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_HOSTNAME          = "unknown"
	TIME_FORMAT_LAYOUT        = "20060102"
	CLOUDWATCH_BUFFER_SIZE    = 200
	CLOUDWATCH_FLUSH_INTERVAL = 15 * time.Second
	LOGFILE_SUFFIX            = "%Y%m%d%H%M.json"
	LOGFILE_MAX_AGE           = time.Duration(86400*40) * time.Second // 40 DAYS
	LOGFILE_ROTATION_TIME     = time.Duration(86400) * time.Second    // 1 DAY
)

func InitLogs(
	activeCloudWatch bool,
	logLevel, groupName, awsRegion, awsProfileName, credentialsFileName, configProfile, streamPrefix string,
	activeFileOutput bool,
	filePath string,
) {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if activeCloudWatch {
		n := time.Now()

		host, err := os.Hostname()

		if err != nil {
			host = DEFAULT_HOSTNAME
		}

		stream := fmt.Sprintf("%s%s_%s_%d", streamPrefix, host, n.Format(TIME_FORMAT_LAYOUT), n.UnixMilli())

		cwHook, err := cloudwatchhook.NewCloudWatchHook(groupName, stream, awsRegion, awsProfileName, CLOUDWATCH_BUFFER_SIZE, CLOUDWATCH_FLUSH_INTERVAL)

		if err != nil {
			log.Println("error initializing logrus:", err)
		}

		logrus.AddHook(cwHook)

	}

	fileName := "app-log"

	if activeFileOutput {
		fileName = path.Join(filePath, fileName)

		writer, err := rotatelogs.New(
			fileName+LOGFILE_SUFFIX,
			rotatelogs.WithLinkName(fileName),
			rotatelogs.WithMaxAge(LOGFILE_MAX_AGE),
			rotatelogs.WithRotationTime(LOGFILE_ROTATION_TIME),
		)

		logrus.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.WarnLevel:  writer,
				logrus.DebugLevel: writer,
				logrus.TraceLevel: writer,
			},
			&logrus.JSONFormatter{},
		))

		if err != nil {
			fmt.Printf("failed to create rotatelogs: %s", err)
		}
	}

	fmt.Printf("Activating logrus logs. activeCloudWatch: %v\n", activeCloudWatch)

	if activeFileOutput || activeCloudWatch {
		logrus.SetOutput(ioutil.Discard)
	}

	switch logLevel {
	case "Error", "Fatal", "Panic":
		logrus.SetLevel(logrus.ErrorLevel)
	case "Warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "Info":
		logrus.SetLevel(logrus.InfoLevel)
	case "Debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "Trace":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

}
