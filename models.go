package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
func (b *Bullet) Deploy() (sdk.DeserializedObject, error) {
	switch b.Method {
	case "GET", "get":
		return b.client.Get(b.URL), nil
	case "POST", "post":
		return b.client.Create(b.Payload), nil
	case "PUT", "put":
		return b.client.Update(b.Payload), nil
	case "DELETE", "delete":
		return b.client.DeleteFromObject(b.Payload), nil
	}
	return sdk.DeserializedObject{}, errors.New("undefined method")
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
	Name      string    `json:"name"`
	Bullets   []Bullet  `json:"bullets"`
	StartTime time.Time `json:"-"`
}

// Drop iterates through the Bullets within a bomb, fires all of them, and
// returns a summary of the results.
func Drop(bomb Bomb) []sdk.DeserializedObject {
	var summary []sdk.DeserializedObject
	for _, bullet := range bomb.Bullets {
		obj, _ := bullet.Deploy()
		summary = append(summary, obj)
	}
	return summary
}

// A Raid is a collection of bombs capable of reporting summary statistics.
type Raid struct {
	StartTime time.Time `json:"-"`
	Bombs     []Bomb    `json:"bombs"`
}

// Conduct iterates through the Bombs in a Raid's Payload, dropping each of
// them, and then returns a summary of the results.
func (r *Raid) Conduct() []byte {
	var raidSummary []byte
	for _, bomb := range r.Bombs {
		response, err := json.MarshalIndent(Drop(bomb), "", "    ")
		check(err)
		log.Debugf("%s", response)
		raidSummary = append(raidSummary, response...)
	}
	return raidSummary
}

// Duration reports how much time has elapsed since the start of the Raid.
func (r *Raid) Duration() time.Duration {
	return time.Now().Sub(r.StartTime)
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
	return Raid{time.Now(), payload}
}

type SimpleBullet struct {
	Name string `json:"name"`
}

type SimpleBomb struct {
	Name    string         `json:"name"`
	Bullets []SimpleBullet `json:"bullets"`
}

type SimpleRaid struct {
	StartTime time.Time    `json:"start_time"`
	Bombs     []SimpleBomb `json:"bombs"`
}
