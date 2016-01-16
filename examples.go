package main

import (
	"time"

	sdk "github.com/dysolution/espsdk"
	"github.com/icrowley/fake"
)

var weapons = make(map[string]Armed)
var arsenals = make(map[string]Arsenal)

func bullet(name string, method string, url string, payload sdk.RESTObject) {
	weapons[name] = Bomb{client, name, method, url, payload}
}

func makeArsenal(name string, weapons ...Armed) {
	var ordnance []Armed
	for _, weapon := range weapons {
		ordnance = append(ordnance, weapon)
	}
	arsenals[name] = Arsenal{
		Name:    name,
		Weapons: ordnance,
	}
}

func foo() Missile {
	return Missile{client, "delete_last_batch", client.DeleteLastBatch}
}

func deleteNewestBatch() Arsenal {
	return Arsenal{
		Name:    "delete_newest_batch",
		Weapons: []Armed{foo()}}
}

func defineBullets() {
	bullet("get_batches", "GET", sdk.Batches, nil)
	bullet("get_a_batch", "GET", "", sdk.Batch{ID: 86102})

	bullet("create_batch", "POST", sdk.Batches, sdk.Batch{
		SubmissionName: appID + ": " + fake.FullName(),
		SubmissionType: "getty_creative_video",
	})

	newBatchData := sdk.Batch{
		ID:             86102,
		SubmissionName: "updated headline",
		Note:           "updated note",
	}
	bullet("update_a_batch", "PUT", "", newBatchData)

	badBatch := sdk.Batch{ID: -1}
	bullet("get_invalid_batch", "GET", badBatch.Path(), badBatch)

	edPhoto := sdk.Contribution{
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
	bullet("create_photo", "POST", edPhoto.Path(), edPhoto)

	edBatch := sdk.Batch{ID: 86103}
	bullet("get_photos", "GET", edBatch.Path(), edBatch)

	release := sdk.Release{
		SubmissionBatchID: 86103,
		FileName:          "some_property.jpg",
		ReleaseType:       "Property",
		FilePath:          "submission/releases/batch_86103/24780225369200015_some_property.jpg",
		MimeType:          "image/jpeg",
	}
	bullet("create_release", "POST", release.Path(), release)
}

// ExampleConfig returns an example of a complete configuration for the app.
// When marshaled into JSON, this can be used as the contents of the config
// file.
func ExampleConfig() Raid {
	defineBullets()

	makeArsenal("batch",
		weapons["create_batch"],
		weapons["get_a_batch"],
		weapons["update_a_batch"],
	)
	makeArsenal("create_batch",
		weapons["create_batch"],
	)
	// makeArsenal("delete_last_batch",
	// 	// weapons["delete_last_batch"],
	// 	getMissile(),
	// )
	// bombs["delete_newest_batch"] = getBomb()
	makeArsenal("get_batch",
		weapons["get_a_batch"],
	)
	makeArsenal("update_batch",
		weapons["update_a_batch"],
	)
	makeArsenal("create_and_confirm_batch",
		weapons["get_batches"],
		weapons["create_batch"],
		weapons["get_batches"],
	)
	makeArsenal("create_and_delete_batch",
		weapons["create_batch"],
		foo(),
	)
	makeArsenal("get_invalid_batches",
		weapons["get_invalid_batch"],
	)
	makeArsenal("create_and_confirm_photo",
		weapons["create_photo"],
		weapons["get_photos"],
	)
	makeArsenal("upload_a_release",
		weapons["create_release"],
	)

	var parallelRaid []Arsenal
	for i := 1; i <= 3; i++ {
		// parallelRaid = append(parallelRaid, bombs["get_batch"])
		parallelRaid = append(parallelRaid, arsenals["create_and_delete_batch"])
	}

	return NewRaid(
		arsenals["create_and_delete_batch"],
		deleteNewestBatch(),
		// parallelRaid...,
	// bombs["batch"],
	// bombs["batch"],
	// bombs["get_batch"],
	// bombs["create_batch"],
	// bombs["create_batch"],
	// bombs["create_batch"],
	// bombs["delete_last_batch"],
	// bombs["create_and_confirm_batch"],
	// bombs["get_invalid_batches"],
	// bombs["create_and_confirm_photo"],
	// bombs["upload_a_release"],
	)
}
