package main

import (
	"time"

	"github.com/dysolution/airstrike"
	"github.com/dysolution/espsdk"
	"github.com/dysolution/sleepwalker"
	"github.com/icrowley/fake"
)

var weapons = make(map[string]airstrike.ArmedWeapon)
var planes = make(map[string]airstrike.Plane)

func makeBomb(name string, method string, url string, payload sleepwalker.RESTObject) {
	weapons[name] = airstrike.Bomb{
		Client:  client,
		Name:    name,
		Method:  method,
		URL:     url,
		Payload: payload,
	}
}

func armPlane(name string, weapons ...airstrike.ArmedWeapon) {
	var ordnance []airstrike.ArmedWeapon
	for _, weapon := range weapons {
		ordnance = append(ordnance, weapon)
	}
	planes[name] = airstrike.Plane{
		Name:    name,
		Arsenal: ordnance,
	}
}

func foo() airstrike.ArmedWeapon {
	return airstrike.Missile{
		Client:    client,
		Name:      "delete_last_batch",
		Operation: espsdk.DeleteLastBatch,
	}
}

func deleteNewestBatch() airstrike.Plane {
	return airstrike.Plane{
		Name:    "delete_newest_batch",
		Arsenal: []airstrike.ArmedWeapon{foo()}}
}

func defineWeapons() {
	makeBomb("get_batches", "GET", espsdk.Batches, nil)
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
func ExampleConfig() airstrike.Raid {
	defineWeapons()

	armPlane("batch",
		weapons["create_batch"],
		weapons["get_a_batch"],
		weapons["update_a_batch"],
	)
	armPlane("create_batch",
		weapons["create_batch"],
	)
	// makeArsenal("delete_last_batch",
	// 	// weapons["delete_last_batch"],
	// 	getMissile(),
	// )
	// bombs["delete_newest_batch"] = getBomb()
	armPlane("get_batch",
		weapons["get_a_batch"],
	)
	armPlane("update_batch",
		weapons["update_a_batch"],
	)
	armPlane("create_and_confirm_batch",
		weapons["get_batches"],
		weapons["create_batch"],
		weapons["get_batches"],
	)
	armPlane("create_and_delete_batch",
		weapons["create_batch"],
		foo(),
	)
	armPlane("get_invalid_batches",
		weapons["get_invalid_batch"],
	)
	armPlane("create_and_confirm_photo",
		weapons["create_photo"],
		weapons["get_photos"],
	)
	armPlane("upload_a_release",
		weapons["create_release"],
	)

	var parallelRaid []airstrike.Plane
	for i := 1; i <= 3; i++ {
		// parallelRaid = append(parallelRaid, bombs["get_batch"])
		parallelRaid = append(parallelRaid, planes["create_and_delete_batch"])
	}

	return airstrike.NewRaid(
		// planes["create_and_delete_batch"],
		// deleteNewestBatch(),
		// parallelRaid...,
		// bombs["batch"],
		// bombs["batch"],
		// bombs["create_batch"],
		// bombs["create_batch"],
		// bombs["create_batch"],
		// bombs["delete_last_batch"],
		// bombs["create_and_confirm_batch"],
		// bombs["get_invalid_batches"],
		// bombs["create_and_confirm_photo"],
		// bombs["upload_a_release"],
		planes["get_batch"],
	)
}
