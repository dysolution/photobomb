// Photobomb conducts workflow tests triggered by requests to its web server.
package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	sdk "github.com/dysolution/espsdk"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

var client = sdk.GetClient(
	os.Getenv("ESP_API_KEY"),
	os.Getenv("ESP_API_SECRET"),
	os.Getenv("ESP_USERNAME"),
	os.Getenv("ESP_PASSWORD"),
	"oregon",
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func batch(w http.ResponseWriter, r *http.Request) {
	goodBatch := sdk.Batch{ID: 86503}
	badBatch := sdk.Batch{ID: -1}
	newBatch := sdk.Batch{SubmissionName: "my batch"}

	raid := NewRaid([]Bomb{
		Bomb{
			Bullet{&client, "GET", goodBatch.Path(), goodBatch},
			Bullet{&client, "GET", badBatch.Path(), badBatch},
			Bullet{&client, "POST", sdk.Batches, newBatch},
		},
	})
	log.Debugf("%v", raid)
	summary := raid.Begin()
	w.Write(summary)
}

func main() {
	http.HandleFunc("/", usage)
	http.HandleFunc("/batch", batch)

	tcpSocket := ":8080"

	log.Infof("listening on %s", tcpSocket)
	http.ListenAndServe(tcpSocket, nil)
}

func usage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("/batch\tdo the batch thing"))
}
