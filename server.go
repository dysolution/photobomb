package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

func runServer() {
	http.HandleFunc("/", status)
	http.HandleFunc("/example", showExampleConfig)
	http.HandleFunc("/config", showConfig)
	http.HandleFunc("/execute", execute)
	// TODO http.HandleFunc("/refresh_token", refreshToken)

	tcpSocket := ":8080"
	log.Infof("listening on %s", tcpSocket)
	http.ListenAndServe(tcpSocket, nil)
}

func status(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"host":   r.RemoteAddr,
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info()

	t, err := template.ParseFiles("index.html")
	check(err)

	routes := make(map[string]string)
	routes["/"] = "display this message"
	routes["/example"] = "display an example config"
	routes["/config"] = "display the current config"
	routes["/execute"] = "execute the current config"

	configJSON, err := json.Marshal(config)
	check(err)

	var simpleConfig SimpleRaid
	err = json.Unmarshal(configJSON, &simpleConfig)
	check(err)

	output, err := json.MarshalIndent(simpleConfig, "", "    ")
	check(err)

	data := struct {
		AppName string
		Routes  map[string]string
		Config  string
		Foo     time.Duration
		Request *http.Request
	}{
		AppName: appID,
		Routes:  routes,
		Config:  string(output),
		Foo:     config.Duration(),
		Request: r,
	}

	err = t.Execute(w, data)
	check(err)
}

func showExampleConfig(w http.ResponseWriter, r *http.Request) {
	output, err := json.MarshalIndent(ExampleConfig(), "", "    ")
	check(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

func showConfig(w http.ResponseWriter, r *http.Request) {
	configJSON, err := json.Marshal(config)
	check(err)

	var simpleConfig SimpleRaid
	err = json.Unmarshal(configJSON, &simpleConfig)
	check(err)

	output, err := json.MarshalIndent(simpleConfig, "", "    ")
	check(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

func execute(w http.ResponseWriter, r *http.Request) {
	raid := config
	summary := raid.Conduct()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(summary)
}
