package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	sdk "github.com/dysolution/espsdk"
)

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

func createBatch(w http.ResponseWriter) sdk.DeserializedObject {
	path := sdk.Batch{ID: 86503}.Path()
	data := sdk.Batch{SubmissionName: "my batch"}
	result := sdk.Create(path, data, &client)
	return result
}

func batch(w http.ResponseWriter, r *http.Request) {
	goodBatch := sdk.Batch{ID: 86503}
	badBatch := sdk.Batch{ID: -1}
	newBatch := sdk.Batch{SubmissionName: "my batch"}

	raid := NewRaid([]Bomb{
		Bomb{
			Flechette{&client, "GET", goodBatch.Path(), goodBatch},
			Flechette{&client, "GET", badBatch.Path(), badBatch},
			Flechette{&client, "POST", sdk.Batches, newBatch},
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
