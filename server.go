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

var inception time.Time
var raidCount = 0
var requestCount = 0
var interval chan int
var toggle = make(chan bool, 1)
var enabled = true

func init() {
	inception = time.Now()
}

func runServer() {
	http.HandleFunc("/", mw(status))
	http.HandleFunc("/attack", mw(attack))
	http.HandleFunc("/config", mw(showConfig))
	http.HandleFunc("/example", mw(showExampleConfig))
	http.HandleFunc("/once", mw(once))
	http.HandleFunc("/pause", mw(pause))
	http.HandleFunc("/warning_shot", mw(once))
	// TODO http.HandleFunc("/refresh_token", refreshToken)

	go func() {
		seconds := 5
		log.Infof("initial interval: %d seconds", seconds)
		for {
			select {
			case seconds := <-interval:
				log.Infof("switching to %d second interval", seconds)
			case enabled = <-toggle:
			default:
			}
			if enabled {
				log.Infof("conducting raid")
				config.Conduct()
				raidCount += 1
				log.Infof("sleeping for %d seconds", seconds)
				time.Sleep(time.Duration(seconds) * time.Second)
			}
		}
	}()

	tcpSocket := ":8080"
	log.Infof("listening on %s", tcpSocket)
	http.ListenAndServe(tcpSocket, nil)
}

func mw(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestCount += 1
		name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		log.WithFields(log.Fields{
			"host":       r.RemoteAddr,
			"method":     r.Method,
			"name":       name,
			"path":       r.URL.Path,
			"request_id": requestCount,
		}).Info()

		fn(w, r)

		log.WithFields(log.Fields{
			"host":          r.RemoteAddr,
			"method":        r.Method,
			"name":          name,
			"path":          r.URL.Path,
			"request_id":    requestCount,
			"response_time": time.Since(start),
		}).Info()
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
	routes["/attack"] = "commence an attack"
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
		Request      *http.Request
		RaidCount    int
		RequestCount int
		Routes       map[string]string
		Uptime       time.Duration
	}{
		AppName:      appID,
		Config:       string(output),
		Enabled:      enabled,
		Request:      r,
		RaidCount:    raidCount,
		RequestCount: requestCount,
		Routes:       routes,
		Uptime:       time.Since(inception),
	}

	err = t.Execute(w, data)
	check(err)
}

func faster(w http.ResponseWriter, r *http.Request) {
	interval <- 4
	w.Write([]byte("faster"))
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
