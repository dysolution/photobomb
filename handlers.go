package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"time"

	"github.com/dysolution/airstrike"
	"github.com/spf13/viper"
)

func mw(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		desc := "httpd"
		requestCount++
		name := "httpd"
		log.WithFields(map[string]interface{}{
			"host":       r.RemoteAddr,
			"method":     r.Method,
			"name":       name,
			"path":       r.URL.Path,
			"request_id": requestCount,
		}).Debug(desc)

		fn(w, r)

		log.WithFields(map[string]interface{}{
			"host":          r.RemoteAddr,
			"method":        r.Method,
			"name":          name,
			"path":          r.URL.Path,
			"request_id":    requestCount,
			"response_time": time.Since(start),
		}).Info(desc)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	fatal(err)

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
	fatal(err)

	var simpleConfig airstrike.SimpleRaid
	err = json.Unmarshal(configJSON, &simpleConfig)
	fatal(err)

	output, err := json.MarshalIndent(simpleConfig, "", "    ")
	fatal(err)

	inception := viper.GetTime("mission.inception")

	data := struct {
		AppName      string
		Config       string
		Enabled      bool
		Interval     float64
		QPS          string
		Request      *http.Request
		RaidCount    int
		RequestCount int
		Routes       map[string]string
		Uptime       time.Duration
	}{
		AppName:      appID,
		Config:       string(output),
		Enabled:      cfg.Mission.Enabled,
		Interval:     cfg.Mission.Interval,
		QPS:          fmt.Sprintf("%.1f", 1000.0/float64(cfg.Mission.Interval)),
		Request:      r,
		RaidCount:    cfg.Mission.RaidCount,
		RequestCount: requestCount,
		Routes:       routes,
		Uptime:       time.Since(inception),
	}

	err = t.Execute(w, data)
	fatal(err)
}

func backoff(w http.ResponseWriter, r *http.Request) {
	newInterval := float64(cfg.Mission.Interval) * math.Phi
	intervalDeltaCh <- float64(newInterval - float64(cfg.Mission.Interval))
	w.Write([]byte(fmt.Sprintf("backing off to %v", newInterval)))
}

func speedup(w http.ResponseWriter, r *http.Request) {
	newInterval := float64(cfg.Mission.Interval) / math.Phi
	intervalDeltaCh <- float64(newInterval - float64(cfg.Mission.Interval))
	w.Write([]byte(fmt.Sprintf("speeding up to %v", newInterval)))
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%s", status)))
}

func pause(w http.ResponseWriter, r *http.Request) {
	cfg.Mission.EnabledCh <- false
	log.Infof("paused")
	w.Write([]byte("paused"))
}

func attack(w http.ResponseWriter, r *http.Request) {
	cfg.Mission.EnabledCh <- true
	log.Infof("attacking")
	w.Write([]byte("attacking"))
}

func showConfig(w http.ResponseWriter, r *http.Request) {
	configJSON, err := json.MarshalIndent(config, "", "  ")
	fatal(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(configJSON)
}

func showConfigNew(w http.ResponseWriter, r *http.Request) {
	configJSON, err := json.MarshalIndent(config, "", "  ")
	fatal(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(configJSON)
}

func showExampleConfig(w http.ResponseWriter, r *http.Request) {
	output, err := json.MarshalIndent(ExampleRaid(), "", "    ")
	fatal(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}
