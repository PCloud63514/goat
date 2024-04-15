package goat

import (
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_LOG_LEVEL string = "debug"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	logLevel := GetPropertyString("LOG_LEVEL", DEFAULT_LOG_LEVEL)
	logLv, err := logrus.ParseLevel(logLevel)
	if nil != err {
		logrus.Warningf("LogLevel=[%v] Parse Failed. default=debug", logLevel)
		logLv = logrus.DebugLevel
	}
	logrus.SetLevel(logLv)
}
