package util

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func InitLogger() {
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	// logrus.SetBufferPool(nil)

	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", ""
			// return " <" + strings.ReplaceAll(frame.Function, "/", ".") + " " + string(frame.Line) + " >", ""
		},
		ForceColors:     true,
		PadLevelText:    true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	pathMap := lfshook.PathMap{
		logrus.InfoLevel: "info.log",
	}
	logrus.AddHook(lfshook.NewHook(
		pathMap,
		&logrus.TextFormatter{},
	))
}
