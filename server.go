package main

import (
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

var inception time.Time
var raidCount, requestCount int
var interval int
var intervalDelta = make(chan float64, 1)
var toggle = make(chan bool, 1)
var enabled bool

var log = logrus.New()

func init() {
	inception = time.Now()
	enabled = false
	interval = 5
	log.Formatter = &prefixed.TextFormatter{TimestampFormat: time.RFC3339}
}

// round returns the nearest integer. This implementation doesn't work for
// negative numbers, but that doesn't matter in this context.
func round(a float64) int {
	val := int(a + 0.5)
	if val < 1 {
		val = 1
	}
	return val
}

func setInterval(d float64) {
	log.Debugf("changing interval by %v seconds", d)
	interval = round(float64(interval) + d)
	log.Debugf("new interval: %v", interval)
}

func httpd() {
	http.HandleFunc("/", mw(status))
	http.HandleFunc("/config", mw(showConfig))
	http.HandleFunc("/example", mw(showExampleConfig))

	http.HandleFunc("/attack", mw(attack))
	http.HandleFunc("/cease_fire", mw(pause))
	http.HandleFunc("/once", mw(once))
	http.HandleFunc("/pause", mw(pause))
	http.HandleFunc("/warning_shot", mw(once))

	http.HandleFunc("/faster", mw(faster))
	http.HandleFunc("/slower", mw(slower))
	http.HandleFunc("/backoff", mw(backoff))
	http.HandleFunc("/speedup", mw(speedup))
	// TODO http.HandleFunc("/refresh_token", refreshToken)

	go func() {
		log.Infof("initial interval: %v seconds", interval)
		for {
			select {
			case enabled = <-toggle:
			case d := <-intervalDelta:
				setInterval(d)
			default:
			}
			if enabled {
				log.Infof("conducting raid")
				config.Conduct()
				raidCount += 1
				log.Infof("sleeping for %v seconds", interval)
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}
	}()

	tcpSocket := ":8080"
	log.Infof("listening on %s", tcpSocket)
	http.ListenAndServe(tcpSocket, nil)
}
