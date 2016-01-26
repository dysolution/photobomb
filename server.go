package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"gopkg.in/fsnotify.v1"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/airstrike"
	"github.com/dysolution/espsdk"
	"github.com/spf13/viper"
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

	go beginMission(&cfg.Mission, reporter)

	tcpSocket := ":8080"
	log.Infof("listening on %s", tcpSocket)
	http.ListenAndServe(tcpSocket, nil)
}

// runs in a goroutine
func beginMission(mission *airstrike.Mission, reporter airstrike.Reporter) {
	desc := "beginMission"

	// set up the reporter for logging and console output
	logFields := make(chan map[string]interface{})
	newThreshold := make(chan time.Duration)

	reporter.ThresholdReceiver = newThreshold

	go reporter.Run(logFields)

	logCh <- map[string]interface{}{
		"severity": "info",
		"source":   desc,
		"interval": mission.Interval,
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		fmt.Println(viper.GetString("foo"))
	})

	for {
		select {
		case mission.Enabled = <-toggle:
		case d := <-intervalDelta:
			setInterval(logFields, d, mission)
			reporter.ThresholdReceiver <- time.Duration(mission.Interval) * time.Millisecond
		default:
		}

		if mission.Enabled {
			logFields <- map[string]interface{}{
				"msg":    "conducting raid",
				"source": desc,
			}

			config.Conduct(log, espsdk.APIInvariant, logFields)
			mission.RaidCount++
			pauseBetweenRaids(logFields, mission)
		}
	}
}

func pauseBetweenRaids(logFields chan map[string]interface{}, mission *airstrike.Mission) {
	desc := "pauseBetweenRaids"
	if mission.Interval <= 0 {
		logFields <- map[string]interface{}{
			"severity": "error",
			"err":      errors.New("non-positive interval; defaulting to 1s"),
			"interval": mission.Interval,
			"source":   desc,
		}
		time.Sleep(1000 * time.Millisecond)
	} else {
		logFields <- map[string]interface{}{
			"msg":      "sleeping",
			"interval": mission.Interval,
			"source":   desc,
		}

		time.Sleep(time.Duration(mission.Interval) * time.Millisecond)

		logFields <- map[string]interface{}{
			"msg":      "waking up",
			"interval": mission.Interval,
			"source":   desc,
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
