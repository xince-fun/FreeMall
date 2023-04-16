package logger

import (
	"github.com/cloudwego/kitex/pkg/klog"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
	"time"
)

// InitKLogger initializes the logger.
func InitKLogger(logDirPath string, logLevel string) {
	// Customize log file path.
	if err := os.MkdirAll(logDirPath, os.ModePerm); err != nil {
		panic(err)
	}

	// Set log filename to date. "2006-01-02" is the format of date.
	logFileName := time.Now().Format("2006-01-02") + ".log"
	filename := path.Join(logDirPath, logFileName)
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
		}
	}
	// 默认的log不打印行号
	logger := kitexlogrus.NewLogger()
	// Provides compression and deletion
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    20,   // A file can be up to 20M.
		MaxBackups: 5,    // Save up to 5 files at the same time.
		MaxAge:     10,   // A file can exist for a maximum of 10 days.
		Compress:   true, // Compress with gzip.
	}

	if runtime.GOOS == "linux" {
		logger.SetOutput(lumberjackLogger)
		logger.SetLevel(klog.LevelWarn)
	} else {
		logger.SetLevel(klog.LevelDebug)
	}

	klog.SetLogger(logger)
}
