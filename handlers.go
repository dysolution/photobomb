package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
)

func mw(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestCount += 1
		name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		log.WithFields(logrus.Fields{
			"host":       r.RemoteAddr,
			"method":     r.Method,
			"name":       name,
			"path":       r.URL.Path,
			"request_id": requestCount,
		}).Debug("httpd")

		fn(w, r)

		log.WithFields(logrus.Fields{
			"host":          r.RemoteAddr,
			"method":        r.Method,
			"name":          name,
			"path":          r.URL.Path,
			"request_id":    requestCount,
			"response_time": time.Since(start),
		}).Info("httpd")
	}
}
func status(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	check(err)

	routes := make(map[string]string)
	routes["/"] = "display this message"
	routes["/example"] = "display an example config"
	routes["/config"] = "display the current config"
	routes["/faster"] = "decrease the interval between attacks"
	routes["/slower"] = "increase the interval between attacks"
	routes["/once"] = "execute the current config once"
	routes["/warning_shot"] = "execute the current config once"
	routes["/attack"] = "commence an attack"
	routes["/speedup"] = "exponentially speed up the attack"
	routes["/backoff"] = "exponentially slow down the attack"
	routes["/cease_fire"] = "pause an attack"
	routes["/pause"] = "pause an attack"

	configJSON, err := json.Marshal(config)
	check(err)

	var simpleConfig SimpleRaid
	err = json.Unmarshal(configJSON, &simpleConfig)
	check(err)

	output, err := json.MarshalIndent(simpleConfig, "", "    ")
	check(err)

	data := struct {
		AppName      string
		Config       string
		Enabled      bool
		Interval     float32
		Request      *http.Request
		RaidCount    int
		RequestCount int
		Routes       map[string]string
		Uptime       time.Duration
	}{
		AppName:      appID,
		Config:       string(output),
		Enabled:      enabled,
		Interval:     interval,
		Request:      r,
		RaidCount:    raidCount,
		RequestCount: requestCount,
		Routes:       routes,
		Uptime:       time.Since(inception),
	}

	err = t.Execute(w, data)
	check(err)
}

func backoff(w http.ResponseWriter, r *http.Request) {
	newInterval := interval * 1.618
	intervalDelta <- newInterval - interval
	w.Write([]byte(fmt.Sprintf("backing off to %v", newInterval)))
}

func speedup(w http.ResponseWriter, r *http.Request) {
	newInterval := interval / 1.618
	intervalDelta <- newInterval - interval
	w.Write([]byte(fmt.Sprintf("speeding up to %v", newInterval)))
}

func faster(w http.ResponseWriter, r *http.Request) {
	intervalDelta <- -1.0
	w.Write([]byte("faster"))
}

func slower(w http.ResponseWriter, r *http.Request) {
	intervalDelta <- 1.0
	w.Write([]byte("slower"))
}

func pause(w http.ResponseWriter, r *http.Request) {
	toggle <- false
	log.Infof("paused")
	w.Write([]byte("paused"))
}

func attack(w http.ResponseWriter, r *http.Request) {
	toggle <- true
	log.Infof("attacking")
	w.Write([]byte("attacking"))
}

func showConfig(w http.ResponseWriter, r *http.Request) {
	configJSON, err := json.MarshalIndent(config, "", "  ")
	check(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(configJSON)
}

func showExampleConfig(w http.ResponseWriter, r *http.Request) {
	output, err := json.MarshalIndent(ExampleConfig(), "", "    ")
	check(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

func once(w http.ResponseWriter, r *http.Request) {
	allResults, err := config.Conduct()
	check(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	output, err := json.MarshalIndent(allResults, "", "  ")
	w.Write(output)
}
