package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/airstrike"
	"github.com/dysolution/sleepwalker"
	"github.com/x-cray/logrus-prefixed-formatter"
)

// NAME appears in the CLI output.
const NAME = "photobomb"

// VERSION appears in the CLI output.
const VERSION = "0.1.0"

var appID = fmt.Sprintf("%s %s", NAME, VERSION)

var client sleepwalker.RESTClient
var cfg Config
var config airstrike.Raid
var intervalDelta = make(chan float64, 1)
var log *logrus.Logger
var logCh = make(chan map[string]interface{})
var logWarning = make(chan map[string]interface{})
var requestCount int
var reporter airstrike.Reporter
var toggle = make(chan bool, 1)
var token sleepwalker.Token
var warningThreshold time.Duration

func init() {
	log = logrus.New()
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

func setInterval(logCh chan map[string]interface{}, d float64, mission *airstrike.Mission) {
	logCh <- map[string]interface{}{
		"message": "changing interval",
		"delta":   d,
	}
	mission.Interval = round(float64(mission.Interval) + d)
	logCh <- map[string]interface{}{
		"interval": mission.Interval,
	}
	if mission.Interval <= 0 {
		logCh <- map[string]interface{}{
			"severity": "error",
			"message":  "invalid interval",
			"interval": mission.Interval,
		}
		mission.Interval = 1
	}
}

func tableFlip(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func loadConfig(path string) airstrike.Raid {
	fi, err := os.Stat(path)
	if err != nil {
		return ExampleConfig()
	}
	if fi.Size() == 0 {
		// The user has probably tried to redirect "example" output, e.g.,
		// photobomb example > config.json, which zeroes out config.json, so
		// we shouldn't bother trying to read it.
		return ExampleConfig()
	}

	log.Debugf("reading configuration from: %s", path)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Debugf("could not read config: %s", path)
		return ExampleConfig()
	}

	var data airstrike.Raid
	if err := json.Unmarshal(file, &data); err != nil {
		log.Fatal(err)
	}
	tableFlip(err)
	return data
}
