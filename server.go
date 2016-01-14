package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"reflect"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
)

const sleepInterval = 5

var inception = time.Now()
var requestCount = 0

func runServer() {
	http.HandleFunc("/", mw(status))
	http.HandleFunc("/attack", mw(attack))
	http.HandleFunc("/config", mw(showConfig))
	http.HandleFunc("/example", mw(showExampleConfig))
	http.HandleFunc("/once", mw(once))
	http.HandleFunc("/warning_shot", mw(once))
	// TODO http.HandleFunc("/refresh_token", refreshToken)

	tcpSocket := ":8080"
	log.Infof("listening on %s", tcpSocket)
	http.ListenAndServe(tcpSocket, nil)
}

func middleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestCount += 1
		name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		log.WithFields(log.Fields{
			"request_id": requestCount,
			"name":       name,
			"host":       r.RemoteAddr,
			"method":     r.Method,
			"path":       r.URL.Path,
		}).Info()
		defer log.WithFields(log.Fields{
			"request_id":    requestCount,
			"name":          name,
			"response_time": time.Since(start),
		}).Info()
		fn(w, r)
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	check(err)

	routes := make(map[string]string)
	routes["/"] = "display this message"
	routes["/example"] = "display an example config"
	routes["/config"] = "display the current config"
	routes["/once"] = "execute the current config once"
	routes["/warning_shot"] = "execute the current config once"
	routes["/attack"] = "execute the current config indefinitely"

	configJSON, err := json.Marshal(config)
	check(err)

	var simpleConfig SimpleRaid
	err = json.Unmarshal(configJSON, &simpleConfig)
	check(err)

	output, err := json.MarshalIndent(simpleConfig, "", "    ")
	check(err)

	data := struct {
		AppName      string
		Routes       map[string]string
		Config       string
		Uptime       time.Duration
		RequestCount int
		Request      *http.Request
	}{
		AppName:      appID,
		Routes:       routes,
		Config:       string(output),
		Uptime:       time.Since(inception),
		RequestCount: requestCount,
		Request:      r,
	}

	err = t.Execute(w, data)
	check(err)
}

func attack(w http.ResponseWriter, r *http.Request) {
	for {
		log.Infof("conducting raid")
		config.Conduct()
		log.Infof("sleeping for %d seconds", sleepInterval)
		time.Sleep(sleepInterval * time.Second)
	}
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
	raid := config
	summary, err := raid.Conduct()
	check(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(summary)
}
