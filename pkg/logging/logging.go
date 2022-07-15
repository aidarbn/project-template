package logging

import (
	logs "github.com/sirupsen/logrus"
)

func Init(logLevel string) {
	// parse string, this is built-in feature of logrus
	ll, err := logs.ParseLevel(logLevel)
	if err != nil {
		ll = logs.DebugLevel
	}
	// set global log level
	logs.SetLevel(ll)
	logs.SetReportCaller(true)
	logs.SetFormatter(&logs.TextFormatter{})
}
