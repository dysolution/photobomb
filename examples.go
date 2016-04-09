package main

import (
	"time"

	as "github.com/dysolution/airstrike"
	"github.com/dysolution/airstrike/ordnance"
	"github.com/dysolution/espsdk"
	sw "github.com/dysolution/sleepwalker"
	"github.com/icrowley/fake"
)

// ExampleConfig returns an example of a complete configuration for the app.
// When marshaled into JSON, this can be used as the contents of the config
// file.
//
// Create Planes that are idempotent, i.e., if you create an object, make
// sure that plane then deletes that object. Otherwise you could end up
// exceeding maximum limits or deleting objects you didn't intend to.
func ExampleRaid() as.Raid {
	armory = ordnance.NewArmory()
	defineWeapons()

	squadron := as.NewSquadron()
	squadron.AddClones(1, client, armory, "get_batch")
	squadron.AddClones(1, client, armory, "get_contributions")

	defineWorkflows(&squadron, client, armory)
	raid, err := as.NewRaid(squadron.Planes...)
	if err != nil {
		return as.Raid{}
	}
	return raid
}

var armory ordnance.Armory

func makeMissile(name string, op func(sw.RESTClient) (sw.Result, error)) {
	armory.NewMissile(client, name, op)
}

func makeBomb(name, method, url string, payload sw.RESTObject) {
	armory.NewBomb(client, name, method, url, payload)
}

// test the Index functionality
func indexes(edImageBatch, crVideoBatch espsdk.Batch) {
	makeBomb("get_batches", "GET", espsdk.BatchesEndpoint, espsdk.Batch{})
	makeMissile(
		"get_contributions",
		func(c sw.RESTClient) (sw.Result, error) {
			// Getting an empty Contribution with a SubmissionBatchID retrieves
			// the index of all Contributions for that Batch.
			return client.Get(espsdk.Contribution{
				SubmissionBatchID: edImageBatch.ID,
			})
		},
	)
	makeMissile(
		"get_releases",
		func(c sw.RESTClient) (sw.Result, error) {
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
	badBatch := espsdk.Batch{ID: "-1"}
	badContribution := espsdk.Contribution{ID: "-1", SubmissionBatchID: edImageBatch.ID}
	badRelease := espsdk.Release{ID: "-1", SubmissionBatchID: crVideoBatch.ID}
	makeBomb("get_invalid_batch", "GET", badBatch.Path(), badBatch)
	makeBomb("get_invalid_contribution", "GET", badContribution.Path(), badContribution)
	makeBomb("get_invalid_release", "GET", badRelease.Path(), badRelease)
}

// test Create / POST functionality
func creates(edImageBatch, crVideoBatch espsdk.Batch) {
	makeBomb("create_batch", "POST", espsdk.BatchesEndpoint, espsdk.Batch{
		SubmissionName:        appID + ": " + fake.FullName(),
		SubmissionType:        "getty_editorial_still",
		SaveExtractedMetadata: false,
	})

	edImage := espsdk.Contribution{
		// reference time: Mon Jan 2 15:04:05 -0700 MST 2006
		CameraShotDate:       time.Now().Format("01/02/2006 15:04:05 -0700"),
		ContentProviderName:  "SHERER, John",
		ContentProviderTitle: "Contributor",
		CountryOfShoot:       "United Kingdom",
		CreditLine:           "John Sherer",
		ExternalFileLocation: "https://c2.staticflickr.com/4/3747/11235643633_60b8701616_o.jpg",
		FileName:             "11235643633_60b8701616_o.jpg",
		Headline:             fake.Sentence(),
		IPTCCategory:         "S",
		SiteDestination:      []string{"Editorial", "WireImage.com"},
		Source:               "AFP",
		SubmissionBatchID:    edImageBatch.ID,
	}
	makeBomb("create_photo", "POST", edImage.Path(), edImage)
	makeMissile("create_and_submit_photo", edImage.CreateAndSubmit)

	release := espsdk.Release{
		FileName:          "some_property.jpg",
		FilePath:          "submission/releases/batch_86572/24780225369200015_some_property.jpg",
		MimeType:          "image/jpeg",
		ReleaseType:       "Property",
		SubmissionBatchID: crVideoBatch.ID,
	}
	makeBomb("create_release", "POST", release.Path(), release)
}

// test Update / PUT functionality
func updates(batch espsdk.Batch, photo espsdk.Contribution) {
	newBatchData := espsdk.Batch{
		ID:             batch.ID,
		Note:           "updated note",
		SubmissionName: "updated headline",
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
func gets(b, c, r sw.RESTObject) {
	makeBomb("get_batch", "GET", "", b)
	makeBomb("get_contribution", "GET", "", c)
	makeBomb("get_release", "GET", "", r)
}

func defineWeapons() {
	// an Editorial Batch and a Creative Batch that are known to exist
	edImageBatch := espsdk.Batch{ID: "86102"}
	crVideoBatch := espsdk.Batch{ID: "89830"}
	contribution := espsdk.Contribution{ID: "1125110", SubmissionBatchID: edImageBatch.ID}
	// contribution := espsdk.Contribution{ID: 1124938, SubmissionBatchID: edImageBatch.ID}
	release := espsdk.Release{ID: "40106", SubmissionBatchID: crVideoBatch.ID}

	creates(edImageBatch, crVideoBatch)
	updates(edImageBatch, contribution)
	indexes(edImageBatch, crVideoBatch)
	gets(edImageBatch, contribution, release)
	duds(edImageBatch, crVideoBatch)

	makeMissile("delete_last_batch", espsdk.DeleteLastBatch)
	makeMissile("submit_photo", contribution.Submit)

}

// A RESTClient can perform operations against a REST API.
type RESTClient interface {
	Get(sw.Findable) (sw.Result, error)
	Create(sw.Findable) (sw.Result, error)
	Update(sw.Findable) (sw.Result, error)
	Delete(sw.Findable) (sw.Result, error)
	Put(sw.Findable, string) (sw.Result, error)
}

func plane(client RESTClient, armory ordnance.Armory, weaponNames ...string) as.Plane {
	var fullName string
	for _, name := range weaponNames {
		fullName += "_" + name
	}
	plane := as.NewPlane(fullName, client)
	arsenal := armory.GetArsenal(weaponNames...)
	plane.Arm(arsenal)
	return plane
}

func defineWorkflows(s *as.Squadron, c RESTClient, a ordnance.Armory) {

	s.Add(plane(c, a,
		"create_batch",
		"delete_last_batch",
	))

	// s.Add(plane(c, a,
	// 	"create_photo",
	// 	"delete_last_photo",
	// ))

	// s.Add(plane(c, a,
	// 	"create_and_submit_photo",
	// ))

	// s.Add(plane(c, a,
	// 	"create_release",
	// 	"delete_last_release",
	// ))

	// s.Add(plane(c, a,
	// 	"get_batch",
	// 	"get_contribution",
	// 	"get_release",
	// ))

	// s.Add(plane(c, a,
	// 	"update_a_batch",
	// 	"update_a_contribution",
	// ))

	// s.Add(plane(c, a,
	// 	"get_batches",
	// 	"get_contributions",
	// 	"get_releases",
	// ))

	// s.Add(plane(c, a,
	// 	"get_invalid_batch",
	// 	"get_invalid_contribution",
	// 	"get_invalid_release",
	// ))
}
