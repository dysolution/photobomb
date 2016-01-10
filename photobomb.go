// Photobomb conducts workflow tests triggered by requests to its web server.
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
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

func runServer() {
	http.HandleFunc("/", usage)
	http.HandleFunc("/example", showExampleConfig)
	http.HandleFunc("/config", showConfig)
	http.HandleFunc("/execute", execute)

	tcpSocket := ":8080"
	log.Infof("listening on %s", tcpSocket)
	http.ListenAndServe(tcpSocket, nil)
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

func showConfig(w http.ResponseWriter, r *http.Request) {
	output, err := json.MarshalIndent(config, "", "    ")
	check(err)
	w.Write(output)
}

func showExampleConfig(w http.ResponseWriter, r *http.Request) {
	output, err := json.MarshalIndent(ExampleConfig(), "", "    ")
	check(err)
	w.Write(output)
}

func execute(w http.ResponseWriter, r *http.Request) {
	raid := config
	summary := raid.Conduct()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(summary)
}

func usage(w http.ResponseWriter, r *http.Request) {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.AppName}}</title>
	</head>
	<body>
    <h1>Routes</h1>
    <table>
      <tr><th>Path</th><th>Action<th></tr>
		{{range $path, $helpText := .Routes}}<tr><td><a href="{{ $path }}">{{ $path }}</a></td><td>{{ $helpText }}</td></tr>{{end}}
    </table>
      <div>
        <h1>Configuration</h2>
        <pre>{{ .Config }}</pre>
      </div>
	</body>
</html>`

	t, err := template.New("webpage").Parse(tpl)
	routes := make(map[string]string)
	routes["/"] = "display this message"
	routes["/example"] = "display an example config"
	routes["/config"] = "display the current config"
	routes["/execute"] = "execute the current config"

	output, err := json.MarshalIndent(config, "", "    ")
	check(err)

	data := struct {
		AppName string
		Routes  map[string]string
		Config  string
	}{
		AppName: appID,
		Routes:  routes,
		Config:  string(output),
	}

	err = t.Execute(w, data)
	check(err)
}
