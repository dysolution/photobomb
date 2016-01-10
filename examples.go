package main

import (
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
		Bullets: bs,
	}
}

func defineBullets() {
	bullet("get_batches", "GET", sdk.Batches, nil)

	bullet("create_batch", "POST", sdk.Batches, sdk.Batch{
		SubmissionName: "photobomb: " + fake.Model(),
		SubmissionType: "getty_creative_video",
	})

	badBatch := sdk.Batch{ID: -1}
	bullet("get_invalid_batch", "GET", badBatch.Path(), badBatch)

	edPhoto := sdk.Contribution{
		SubmissionBatchID:    86102,
		CameraShotDate:       "12/14/2015",
		ContentProviderName:  "provider",
		ContentProviderTitle: "Contributor",
		CountryOfShoot:       "United States",
		CreditLine:           "Ansel Adams",
		FileName:             "el_capitan_merced_river_clouds.jpg",
		Headline:             "El Capitan, Merced River, Clouds",
		IptcCategory:         "S",
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

	makeBomb("create_and_confirm_batch",
		bullets["get_batches"],
		bullets["create_batch"],
		bullets["get_batches"],
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

	return NewRaid(
		bombs["create_and_confirm_batch"],
		bombs["get_invalid_batches"],
		bombs["create_and_confirm_photo"],
		bombs["upload_a_release"],
	)
}
