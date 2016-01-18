package main

import (
	"time"

	as "github.com/dysolution/airstrike"
	"github.com/dysolution/airstrike/arsenal"
	"github.com/dysolution/espsdk"
	"github.com/dysolution/sleepwalker"
	"github.com/icrowley/fake"
)

var w = make(map[string]arsenal.ArmedWeapon)
var a = make(map[string]arsenal.Arsenal)
var p = make(map[string]as.Plane)

var planes []as.Plane

func makeBomb(name string, method string, url string, payload sleepwalker.RESTObject) {
	w[name] = as.Bomb{
		Client:  client,
		Name:    name,
		Method:  method,
		URL:     url,
		Payload: payload,
	}
}

func deleteNewestBatch() arsenal.ArmedWeapon {
	return as.Missile{
		Client:    client,
		Name:      "delete_last_batch",
		Operation: espsdk.DeleteLastBatch,
	}
}

func defineWeapons() {
	makeBomb("get_batches", "GET", espsdk.Batches, espsdk.Batch{})
	makeBomb("get_a_batch", "GET", "", espsdk.Batch{ID: 86102})

	makeBomb("create_batch", "POST", espsdk.Batches, espsdk.Batch{
		SubmissionName: appID + ": " + fake.FullName(),
		SubmissionType: "getty_creative_video",
	})

	newBatchData := espsdk.Batch{
		ID:             86102,
		SubmissionName: "updated headline",
		Note:           "updated note",
	}
	makeBomb("update_a_batch", "PUT", "", newBatchData)

	badBatch := espsdk.Batch{ID: -1}
	makeBomb("get_invalid_batch", "GET", badBatch.Path(), badBatch)

	edPhoto := espsdk.Contribution{
		SubmissionBatchID:    86102,
		CameraShotDate:       time.Now().Format("01/02/2006"),
		ContentProviderName:  "provider",
		ContentProviderTitle: "Contributor",
		CountryOfShoot:       fake.Country(),
		CreditLine:           fake.FullName(),
		FileName:             fake.Word() + ".jpg",
		Headline:             fake.Sentence(),
		IPTCCategory:         "S",
		SiteDestination:      []string{"Editorial", "WireImage.com"},
		Source:               "AFP",
	}
	makeBomb("create_photo", "POST", edPhoto.Path(), edPhoto)

	edBatch := espsdk.Batch{ID: 86103}
	makeBomb("get_photos", "GET", edBatch.Path(), edBatch)

	release := espsdk.Release{
		SubmissionBatchID: 86103,
		FileName:          "some_property.jpg",
		ReleaseType:       "Property",
		FilePath:          "submission/releases/batch_86103/24780225369200015_some_property.jpg",
		MimeType:          "image/jpeg",
	}
	makeBomb("create_release", "POST", release.Path(), release)
}

// ExampleConfig returns an example of a complete configuration for the app.
// When marshaled into JSON, this can be used as the contents of the config
// file.
func ExampleConfig() as.Raid {
	defineWeapons()

	a["batch"] = arsenal.New(w["create_batch"], w["get_a_batch"], w["update_a_batch"])
	a["create_and_confirm_batch"] = arsenal.New(w["get_batches"], w["create_batch"], w["get_batches"])
	a["create_and_confirm_photo"] = arsenal.New(w["create_photo"], w["get_photos"])
	a["create_and_delete_batch"] = arsenal.New(w["create_batch"], deleteNewestBatch())
	a["create_batch"] = arsenal.New(w["create_batch"])
	a["create_batch"] = arsenal.New(w["create_batch"])
	a["create_relearsenale"] = arsenal.New(w["create_relearsenale"])
	a["delete_newest_batch"] = arsenal.New(deleteNewestBatch())
	a["get_batch"] = arsenal.New(w["get_a_batch"])
	a["get_invalid_batches"] = arsenal.New(w["get_invalid_batch"])
	a["update_batch"] = arsenal.New(w["update_a_batch"])

	for name, arsenal := range a {
		plane := as.NewPlane(name, client)
		err := plane.Arm(arsenal)
		if err != nil {
			log.Error(err)
		}
		p[name] = plane
	}

	log.Debug(p)
	return getRaid(30, planes)
}

func getRaid(planeCount int, arsenals []as.Plane) as.Raid {
	var mission []as.Plane
	for i := 1; i <= planeCount; i++ {
		// mission = append(mission, p["create_and_confirm_batch"])
		mission = append(mission, p["create_and_confirm_photo"])
		mission = append(mission, p["create_and_delete_batch"])
		mission = append(mission, p["create_release"])
		mission = append(mission, p["get_batch"])
		// planes["batch"],
		// planes["batch"],
		// planes["create_batch"],
		// planes["delete_last_batch"],
		// planes["get_invalid_batches"],
		// planes["create_batch"],
		// planes["get_batch"],
	}
	return as.NewRaid(mission...)
}
