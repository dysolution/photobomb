package main

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/airstrike"
	"github.com/dysolution/airstrike/ordnance"
	"github.com/dysolution/espsdk"
	"github.com/dysolution/sleepwalker"
	"github.com/icrowley/fake"
)

var deprecatedPlanes []airstrike.Plane

var armory ordnance.Armory

func makeMissile(name string, op func(sleepwalker.RESTClient) (sleepwalker.Result, error)) {
	armory.NewMissile(client, name, op)
}

func makeBomb(name, method, url string, payload sleepwalker.RESTObject) {
	armory.NewBomb(client, name, method, url, payload)
}

// test the Index functionality
func indexes(edImageBatch, crVideoBatch espsdk.Batch) {
	makeBomb("get_batches", "GET", espsdk.Batches, espsdk.Batch{})
	makeMissile(
		"get_contributions",
		func(c sleepwalker.RESTClient) (sleepwalker.Result, error) {
			// Getting an empty Contribution with a SubmissionBatchID retrieves
			// the index of all Contributions for that Batch.
			return client.Get(espsdk.Contribution{
				SubmissionBatchID: edImageBatch.ID,
			})
		},
	)
	makeMissile(
		"get_releases",
		func(c sleepwalker.RESTClient) (sleepwalker.Result, error) {
			// Getting an empty Release with a SubmissionBatchID retrieves
			// the index of all Releases for that Batch.
			return client.Get(espsdk.Release{
				SubmissionBatchID: crVideoBatch.ID,
			})
		},
	)
}

// a few objects that don't exist, to test 404
func duds(edImageBatch, crVideoBatch espsdk.Batch) {
	badBatch := espsdk.Batch{ID: -1}
	badContribution := espsdk.Contribution{ID: -1, SubmissionBatchID: edImageBatch.ID}
	badRelease := espsdk.Release{ID: -1, SubmissionBatchID: crVideoBatch.ID}
	makeBomb("get_invalid_batch", "GET", badBatch.Path(), badBatch)
	makeBomb("get_invalid_contribution", "GET", badContribution.Path(), badContribution)
	makeBomb("get_invalid_release", "GET", badRelease.Path(), badRelease)
}

// test Create / POST functionality
func creates(edImageBatch, crVideoBatch espsdk.Batch) {
	makeBomb("create_batch", "POST", espsdk.Batches, espsdk.Batch{
		SubmissionName: appID + ": " + fake.FullName(),
		SubmissionType: "getty_creative_video",
	})

	edImage := espsdk.Contribution{
		SubmissionBatchID:    edImageBatch.ID,
		CameraShotDate:       time.Now().Format("01/02/2006"),
		ContentProviderName:  "provider",
		ContentProviderTitle: "Contributor",
		CountryOfShoot:       "United Kingdom",
		CreditLine:           "John Sherer",
		ExternalFileLocation: "https://c2.staticflickr.com/4/3747/11235643633_60b8701616_o.jpg",
		FileName:             "11235643633_60b8701616_o.jpg",
		Headline:             fake.Sentence(),
		IPTCCategory:         "S",
		SiteDestination:      []string{"Editorial", "WireImage.com"},
		Source:               "AFP",
	}
	makeBomb("create_photo", "POST", edImage.Path(), edImage)

	release := espsdk.Release{
		SubmissionBatchID: crVideoBatch.ID,
		FileName:          "some_property.jpg",
		ReleaseType:       "Property",
		FilePath:          "submission/releases/batch_86572/24780225369200015_some_property.jpg",
		MimeType:          "image/jpeg",
	}
	makeBomb("create_release", "POST", release.Path(), release)
}

// test Update / PUT functionality
func updates(batch espsdk.Batch, photo espsdk.Contribution) {
	newBatchData := espsdk.Batch{
		ID:             batch.ID,
		SubmissionName: "updated headline",
		Note:           "updated note",
	}
	makeBomb("update_a_batch", "PUT", "", newBatchData)

	contributionUpdate := espsdk.Contribution{
		ID:                photo.ID,
		SubmissionBatchID: batch.ID,
		Headline:          fake.Sentence(),
	}
	makeBomb("update_a_contribution", "PUT", "", contributionUpdate)
}

// test Get / GET functionality
func gets(b espsdk.Batch, c espsdk.Contribution, r espsdk.Release) {
	makeBomb("get_batch", "GET", "", b)
	makeBomb("get_contribution", "GET", "", c)
	makeBomb("get_release", "GET", "", r)
}

func defineWeapons() {
	// an Editorial Batch and a Creative Batch that are known to exist
	edImageBatch := espsdk.Batch{ID: 86102}
	crVideoBatch := espsdk.Batch{ID: 88086}
	contribution := espsdk.Contribution{ID: 1124654, SubmissionBatchID: edImageBatch.ID}
	release := espsdk.Release{ID: 39969, SubmissionBatchID: crVideoBatch.ID}

	creates(edImageBatch, crVideoBatch)
	updates(edImageBatch, contribution)
	indexes(edImageBatch, crVideoBatch)
	gets(edImageBatch, contribution, release)
	duds(edImageBatch, crVideoBatch)

	makeMissile("delete_last_batch", espsdk.DeleteLastBatch)

}

// ExampleConfig returns an example of a complete configuration for the app.
// When marshaled into JSON, this can be used as the contents of the config
// file.
//
// Create Planes that are idempotent, i.e., if you create an object, make
// sure that plane then deletes that object. Otherwise you could end up
// exceeding maximum limits or deleting objects you didn't intend to.
func ExampleConfig() airstrike.Raid {

	armory = ordnance.NewArmory(log)
	defineWeapons()

	squadron := airstrike.New(log)

	// // Krieger makes everything possible.
	// Krieger := airstrike.NewPlane("Krieger", client)
	// Krieger.Arm(armory.GetArsenal(
	// 	"create_batch",
	// 	"delete_last_batch",
	// ))
	// squadron.Add(Krieger)

	// // Cheryl's working.
	// plane = airstrike.NewPlane("Cheryl", client)
	// plane.Arm(armory.GetArsenal(
	// 	"create_photo",
	// 	"delete_last_photo",
	// ))
	// squadron.Add(plane)

	// // // Pam needs you to fill out this form.
	// Pam := airstrike.NewPlane("Pam", client)
	// Pam.Arm(armory.GetArsenal(
	// 	"create_release",
	// 	"delete_last_release",
	// ))
	// squadron.Add(Pam)

	// // Ray doesn't want to hear too much information.
	// Ray := airstrike.NewPlane("Ray", client)
	// Ray.Arm(armory.GetArsenal(
	// 	"get_batch",
	// 	"get_contribution",
	// 	"get_release",
	// ))
	// squadron.Add(Ray)

	// // Archer wants things his way.
	// Archer := airstrike.NewPlane("Archer", client)
	// Archer.Arm(armory.GetArsenal(
	// 	"update_a_batch",
	// 	"update_a_contribution",
	// ))
	// squadron.Add(Archer)

	// // Cyril accounts for everything.
	// Cyril := airstrike.NewPlane("Cyril", client)
	// Cyril.Arm(armory.GetArsenal(
	// 	"get_batches",
	// 	"get_contributions",
	// 	"get_releases",
	// ))
	// squadron.Add(Cyril)

	// // Malory makes unreasonable demands.
	// Malory := airstrike.NewPlane("Malory", client)
	// Malory.Arm(armory.GetArsenal(
	// 	"get_invalid_batch",
	// 	"get_invalid_contribution",
	// 	"get_invalid_release",
	// ))
	// squadron.Add(Malory)

	// You can also simulate heavy load by creating many anonymous Planes
	// that each performs any workflow composed of a single operation or many.
	//
	squadron.AddClones(5, client, armory, "get_batch")
	// AddChaos(5, 2, &squadron)

	return airstrike.NewRaid(squadron.Planes...)
}

// AddChaos creates the specified number of Planes that each have their own
// random selection of n weapons, quantified by their "deadliness."
func AddChaos(clones int, deadliness int, squadron *airstrike.Squadron) {
	desc := "AddChaos"
	for i := 1; i <= clones; i++ {
		name := fmt.Sprintf("chaos_%d_of_%d", i, clones)
		plane := airstrike.NewPlane(name, client)

		weaponNames := armory.GetRandomWeaponNames(deadliness)
		var arsenal ordnance.Arsenal
		for _, weaponName := range weaponNames {
			log.WithFields(logrus.Fields{
				"plane":  name,
				"weapon": weaponName,
				"msg":    "adding weapon",
			}).Debug(desc)
			arsenal = append(arsenal, armory.GetWeapon(weaponName))
		}
		plane.Arm(arsenal)
		squadron.Add(plane)
	}
}
