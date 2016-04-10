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

var (
	appID            = fmt.Sprintf("%s %s", NAME, VERSION)
	cfg              Config
	client           sleepwalker.RESTClient
	config           airstrike.Raid
	log              *logrus.Logger
	reporter         *airstrike.Reporter
	requestCount     int
	status           int
	statusCh         = make(chan int, 1)
	toggle           = make(chan bool, 1)
	token            sleepwalker.Token
	warningThreshold time.Duration
)

type Fields map[string]interface{}

func init() {
	log = logrus.New()
	log.Formatter = &prefixed.TextFormatter{TimestampFormat: time.RFC3339}
}

func fatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func loadConfig(path string) airstrike.Raid {
	fi, err := os.Stat(path)
	if err != nil {
		return ExampleRaid()
	}
	if fi.Size() == 0 {
		// The user has probably tried to redirect "example" output, e.g.,
		// photobomb example > config.json, which zeroes out config.json, so
		// we shouldn't bother trying to read it.
		return ExampleRaid()
	}

	log.Debugf("reading configuration from: %s", path)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Debugf("could not read config: %s", path)
		return ExampleRaid()
	}

	var data airstrike.Raid
	if err := json.Unmarshal(file, &data); err != nil {
		log.Fatal(err)
	}
	fatal(err)
	return data
}
