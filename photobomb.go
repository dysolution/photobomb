// Photobomb conducts workflow tests triggered by requests to its web server.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	sdk "github.com/dysolution/espsdk"
)

const NAME = "photobomb"
const VERSION = "0.0.1"

var appID = fmt.Sprintf("%s %s", NAME, VERSION)

var client sdk.Client
var token sdk.Token
var config Raid

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func loadConfig(path string) Raid {
	fi, err := os.Stat(path)
	check(err)
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

	var data Raid
	if err := json.Unmarshal(file, &data); err != nil {
		log.Fatal(err)
	}
	check(err)
	return data
}
