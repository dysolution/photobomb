package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/airstrike"
	"github.com/dysolution/espsdk"
)

func serve() {
	http.HandleFunc("/", mw(status))
	http.HandleFunc("/config", mw(showConfig))
	http.HandleFunc("/example", mw(showExampleConfig))

	http.HandleFunc("/attack", mw(attack))
	http.HandleFunc("/cease_fire", mw(pause))
	http.HandleFunc("/pause", mw(pause))

	http.HandleFunc("/faster", mw(faster))
	http.HandleFunc("/slower", mw(slower))
	http.HandleFunc("/backoff", mw(backoff))
	http.HandleFunc("/speedup", mw(speedup))
	// TODO http.HandleFunc("/refresh_token", refreshToken)

	go beginMission(reporter)

	tcpSocket := ":8080"
	log.Infof("listening on %s", tcpSocket)
	http.ListenAndServe(tcpSocket, nil)
}

// runs in a goroutine
func beginMission(reporter airstrike.Reporter) {
	desc := "beginMission"

	// set up the reporter for logging and console output
	logFields := make(chan map[string]interface{})

	go reporter.Run(logFields)

	logCh <- map[string]interface{}{
		"severity": "info",
		"source":   desc,
		"interval": interval,
	}

	for {
		select {
		case enabled = <-toggle:
		case d := <-intervalDelta:
			setInterval(d)
		default:
		}
		if enabled {
			logFields <- map[string]interface{}{
				"msg":    "conducting raid",
				"source": desc,
			}

			config.Conduct(log, espsdk.APIInvariant, logFields)
			raidCount++

			logFields <- map[string]interface{}{
				"msg":      "sleeping",
				"interval": interval,
				"source":   desc,
			}
			time.Sleep(time.Duration(interval) * time.Second)
			logFields <- map[string]interface{}{
				"msg":      "waking up",
				"interval": interval,
				"source":   desc,
			}
		}
	}
}

func logger(ch chan map[string]interface{}, log *logrus.Logger) {
	for {
		select {
		case fields := <-ch:
			desc := fmt.Sprintf("photobomb.logger (gr: %v)", runtime.NumGoroutine())
			switch fields["severity"] {
			case "INFO", "info":
				log.WithFields(fields).Info(desc)
			case "WARN", "warn":
				log.WithFields(fields).Warn(desc)
			default:
				log.WithFields(fields).Debug(desc)
			}
		default:
		}
	}
}
