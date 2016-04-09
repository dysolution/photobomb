package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

type Reporter struct {
	Logger *logrus.Logger
	LogCh  chan map[string]interface{}
}

func NewReporter() *Reporter {
	return &Reporter{
		Logger: logrus.New(),
		LogCh:  make(chan map[string]interface{}, 1),
	}
}

func (r *Reporter) listen() {
	desc := fmt.Sprintf("photobomb.logger")
	for {
		select {
		case fields := <-r.LogCh:
			switch fields["severity"] {
			case "INFO", "info":
				r.Logger.WithFields(fields).Info(desc)
			case "WARN", "warn":
				r.Logger.WithFields(fields).Warn(desc)
			default:
				r.Logger.WithFields(fields).Debug(desc)
			}
		default:
		}
	}
}
