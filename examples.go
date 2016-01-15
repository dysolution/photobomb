package main

import (
	"time"

	sdk "github.com/dysolution/espsdk"
	"github.com/icrowley/fake"
)

var bullets = make(map[string]Bullet)
var bombs = make(map[string]Bomb)

func bullet(name string, method string, url string, payload sdk.RESTObject) {
	bullets[name] = Bullet{&client, name, method, url, payload}
}

func makeBomb(name string, bullets ...Bullet) {
	var bs []Bullet
	for _, bullet := range bullets {
		bs = append(bs, bullet)
	}
	bombs[name] = Bomb{
		Name:    name,
		Bullets: bs,
	}
}

func lastBatch() sdk.Batch {
	return client.Get(sdk.Batches).Last()
}

func defineBullets() {
	bullet("get_batches", "GET", sdk.Batches, nil)
	bullet("get_a_batch", "GET", "", sdk.Batch{ID: 86102})

	bullet("create_batch", "POST", sdk.Batches, sdk.Batch{
		SubmissionName: appID + ": " + fake.FullName(),
		SubmissionType: "getty_creative_video",
	})

	bullet("delete_last_batch", "DELETE", "", lastBatch())

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

	makeBomb("batch",
		bullets["create_batch"],
		bullets["get_a_batch"],
		bullets["update_a_batch"],
	)
	makeBomb("create_batch",
		bullets["create_batch"],
	)
	makeBomb("delete_last_batch",
		bullets["delete_last_batch"],
	)
	makeBomb("get_batch",
		bullets["get_a_batch"],
	)
	makeBomb("update_batch",
		bullets["update_a_batch"],
	)
	makeBomb("create_and_confirm_batch",
		bullets["get_batches"],
		bullets["create_batch"],
		bullets["get_batches"],
	)
	makeBomb("create_and_delete_batch",
		bullets["create_batch"],
		bullets["delete_last_batch"],
	)
	makeBomb("get_invalid_batches",
		bullets["get_invalid_batch"],
	)
	makeBomb("create_and_confirm_photo",
		bullets["create_photo"],
		bullets["get_photos"],
	)
	makeBomb("upload_a_release",
		bullets["create_release"],
	)

	var smallRaid []Bomb
	for i := 1; i <= 2; i++ {
		smallRaid = append(smallRaid, bombs["get_batch"])
	}

	return NewRaid(
		smallRaid...,
	// bombs["batch"],
	// bombs["batch"],
	// bombs["get_batch"],
	// bombs["create_batch"],
	// bombs["delete_last_batch"],
	// bombs["create_and_confirm_batch"],
	// bombs["create_and_delete_batch"],
	// bombs["get_invalid_batches"],
	// bombs["create_and_confirm_photo"],
	// bombs["upload_a_release"],
	)
}
