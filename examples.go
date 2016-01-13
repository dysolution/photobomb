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

	bullet("create_batch", "POST", sdk.Batches, sdk.Batch{
		SubmissionName: appID + ": " + fake.FullName(),
		SubmissionType: "getty_creative_video",
	})

	bullet("delete_last_batch", "DELETE", lastBatch().Path(), nil)

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

	makeBomb("create_batch",
		bullets["create_batch"],
	)
	makeBomb("delete_last_batch",
		bullets["delete_last_batch"],
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

	return NewRaid(
		bombs["create_batch"],
		bombs["delete_last_batch"],
	)
}
