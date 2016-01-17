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

func armPlane(name string, client sleepwalker.RESTClient, weapons ...airstrike.ArmedWeapon) {
	var ordnance []airstrike.ArmedWeapon
	for _, weapon := range weapons {
		ordnance = append(ordnance, weapon)
	}
	planes[name] = airstrike.Plane{
		Name:    name,
		Client:  client,
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
func ExampleConfig() airstrike.Raid {
	defineWeapons()

	armPlane("batch",
		client,
		weapons["create_batch"],
		weapons["get_a_batch"],
		weapons["update_a_batch"],
	)
	armPlane("create_batch",
		client,
		weapons["create_batch"],
	)
	// makeArsenal("delete_last_batch",
	// 	// weapons["delete_last_batch"],
	// 	getMissile(),
	// )
	// bombs["delete_newest_batch"] = getBomb()
	armPlane("get_batch",
		client,
		weapons["get_a_batch"],
	)
	armPlane("update_batch",
		client,
		weapons["update_a_batch"],
	)
	armPlane("create_and_confirm_batch",
		client,
		weapons["get_batches"],
		weapons["create_batch"],
		weapons["get_batches"],
	)
	armPlane("create_and_delete_batch",
		client,
		weapons["create_batch"],
		foo(),
	)
	armPlane("get_invalid_batches",
		client,
		weapons["get_invalid_batch"],
	)
	armPlane("create_and_confirm_photo",
		client,
		weapons["create_photo"],
		weapons["get_photos"],
	)
	armPlane("upload_a_release",
		client,
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
		// planes["batch"],
		// planes["batch"],
		// planes["create_batch"],
		// planes["delete_last_batch"],
		// planes["get_invalid_batches"],
		planes["create_and_confirm_batch"],
		planes["create_and_confirm_photo"],
		planes["upload_a_release"],
		// planes["create_batch"],
		// planes["get_batch"],
	)
}
