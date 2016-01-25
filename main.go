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
var config airstrike.Raid
var enabled bool
var interval int
var intervalDelta = make(chan float64, 1)
var log *logrus.Logger
var logCh = make(chan map[string]interface{})
var logWarning = make(chan map[string]interface{})
var raidCount, requestCount int
var reporter airstrike.Reporter
var toggle = make(chan bool, 1)
var token sleepwalker.Token
var warningThreshold time.Duration

func init() {
	enabled = true
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

func setInterval(d float64) {
	log.Debugf("changing interval by %v seconds", d)
	interval = round(float64(interval) + d)
	log.Debugf("new interval: %v", interval)
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
