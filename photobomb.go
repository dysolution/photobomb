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

// foo creates a batch, uploads some photos, deletes some photos,
// and lists the photos in the batch
func foo(w http.ResponseWriter, r *http.Request) {
	goodBatch := sdk.Batch{ID: 86503}
	badBatch := sdk.Batch{ID: -1}
	newBatch := sdk.Batch{SubmissionName: "my batch"}
	raid := NewRaid([]Bomb{
		Bomb{
			Flechette{client, "GET", goodBatch.Path(), goodBatch},
			Flechette{client, "GET", badBatch.Path(), badBatch},
			Flechette{client, "POST", sdk.Batches, newBatch},
		},
	})
	summary := raid.Begin()
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// json.NewEncoder(w).Encode(summary)
	w.Write(summary)
	//createBatch(w)

	// uploadPhotos()
	// deletePhotos()
	// listPhotos()
}

func main() {
	http.HandleFunc("/", usage)
	http.HandleFunc("/batch", foo)

	http.ListenAndServe(":8080", nil)
}

func usage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("/batch\tdo the batch thing"))
}
