package main

import (
	"time"

	sdk "github.com/dysolution/espsdk"
	"github.com/icrowley/fake"
)

var weapons = make(map[string]Armed)
var planes = make(map[string]Arsenal)

func makeBomb(name string, method string, url string, payload sdk.RESTObject) {
	weapons[name] = Bomb{client, name, method, url, payload}
}

func armPlane(name string, weapons ...Armed) {
	var ordnance []Armed
	for _, weapon := range weapons {
		ordnance = append(ordnance, weapon)
	}
	planes[name] = Arsenal{
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

func defineWeapons() {
	makeBomb("get_batches", "GET", sdk.Batches, nil)
	makeBomb("get_a_batch", "GET", "", sdk.Batch{ID: 86102})

	makeBomb("create_batch", "POST", sdk.Batches, sdk.Batch{
		SubmissionName: appID + ": " + fake.FullName(),
		SubmissionType: "getty_creative_video",
	})

	newBatchData := sdk.Batch{
		ID:             86102,
		SubmissionName: "updated headline",
		Note:           "updated note",
	}
	makeBomb("update_a_batch", "PUT", "", newBatchData)

	badBatch := sdk.Batch{ID: -1}
	makeBomb("get_invalid_batch", "GET", badBatch.Path(), badBatch)

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
	makeBomb("create_photo", "POST", edPhoto.Path(), edPhoto)

	edBatch := sdk.Batch{ID: 86103}
	makeBomb("get_photos", "GET", edBatch.Path(), edBatch)

	release := sdk.Release{
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
func ExampleConfig() Raid {
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

	var parallelRaid []Arsenal
	for i := 1; i <= 3; i++ {
		// parallelRaid = append(parallelRaid, bombs["get_batch"])
		parallelRaid = append(parallelRaid, planes["create_and_delete_batch"])
	}

	return NewRaid(
		planes["create_and_delete_batch"],
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
