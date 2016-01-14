package main

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	sdk "github.com/dysolution/espsdk"
)

type Serializable interface {
	Marshal() ([]byte, error)
}

// A Bullet represents a single HTTP request that performs an operation
// against a single API endpoint. Each Bomb can contain one or multiple
// Bullets.
type Bullet struct {
	client  *sdk.Client
	Name    string         `json:"name"`
	Method  string         `json:"method"`
	URL     string         `json:"url"`
	Payload sdk.RESTObject `json:"payload,omitempty"`
}

// Deploy sets the Bullet in motion.
func (b *Bullet) Deploy() (*sdk.FulfilledRequest, error) {
	switch b.Method {
	case "GET", "get":
		fRequest, err := b.client.VerboseGet(b.Payload)
		if err != nil {
			log.Errorf("%s.Deploy: %v", b.Name, err)
			return &sdk.FulfilledRequest{}, err
		}
		log.WithFields(fRequest.Stats()).Infof("%s.Deploy", b.Name)
		return fRequest, nil
	case "POST", "post":
		b.client.Create(b.Payload)
		return &sdk.FulfilledRequest{}, nil
	case "PUT", "put":
		b.client.Update(b.Payload)
		return &sdk.FulfilledRequest{}, nil
	case "DELETE", "delete":
		b.client.DeleteFromObject(b.Payload)
		return &sdk.FulfilledRequest{}, nil
	}
	return &sdk.FulfilledRequest{}, errors.New("undefined method")
}

func (b *Bullet) String() string {
	out, err := json.MarshalIndent(b, "", "  ")
	check(err)
	return fmt.Sprintf("%s", out)
}

// A Bomb is a collection of Bullets. It represents a list of tasks that
// compose a workflow.
//
// For example, a common workflow would be:
//   1. list all batches
//   2. get the metadata for a batch
//   3. upload a contribution to the batch
type Bomb struct {
	Name    string   `json:"name"`
	Bullets []Bullet `json:"bullets"`
}

// Drop iterates through the Bullets within a bomb, fires all of them, and
// returns a summary of the results.
func Drop(bomb Bomb) []*sdk.FulfilledRequest {
	var summary []*sdk.FulfilledRequest
	for _, bullet := range bomb.Bullets {
		result, err := bullet.Deploy()
		if err != nil {
			log.Errorf("Drop: %v", err)
		}
		summary = append(summary, result)
	}
	return summary
}

// A Raid is a collection of bombs capable of reporting summary statistics.
type Raid struct {
	Bombs []Bomb `json:"bombs"`
}

// Conduct iterates through the Bombs in a Raid's Payload, dropping each of
// them, and then returns a summary of the results.
func (r *Raid) Conduct() ([]byte, error) {
	var raidSummary []byte
	for _, bomb := range r.Bombs {
		response, err := json.MarshalIndent(Drop(bomb), "", "    ")
		if err != nil {
			log.Errorf("Raid.Conduct(): %s", err)
			return []byte{}, err
		}
		log.Debugf("%s", response)
		raidSummary = append(raidSummary, response...)
	}
	return raidSummary, nil
}

func (r *Raid) String() string {
	out, err := json.MarshalIndent(r, "", "  ")
	check(err)
	return string(out)
}

// NewRaid initializes and returns a Raid, . It should be used in lieu of Raid literals.
func NewRaid(bombs ...Bomb) Raid {
	var payload []Bomb
	for _, bomb := range bombs {
		payload = append(payload, bomb)
	}
	return Raid{Bombs: payload}
}

type SimpleBullet struct {
	Name string `json:"name"`
}

type SimpleBomb struct {
	Name    string         `json:"name"`
	Bullets []SimpleBullet `json:"bullets"`
}

type SimpleRaid struct {
	Bombs []SimpleBomb `json:"bombs"`
}
