package logger

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	if Log != nil {
		fileName := getFileDir()
		writter := rotateLog(fileName)
		Log.Out = writter
		return
	}

	logger := logrus.New()
	fileName := getFileDir()
	writter := rotateLog(fileName)
	logger.Out = writter
	logger.SetLevel(logrus.DebugLevel)
	Log = logger
}

func getFileDir() string {
	now := time.Now()
	_, filePath, _, _ := runtime.Caller(0)

	// ../../../logs
	logsPath := filepath.Join(filePath, "..", "..", "..", "logs")

	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logsPath, logFileName)

	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
		}
	}

	return fileName
}

func rotateLog(fileName string) *rotatelogs.RotateLogs {
	writter, _ := rotatelogs.New(
		fileName+"%H%M",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(time.Duration(12)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(3)*time.Hour),
	)
	return writter
}
