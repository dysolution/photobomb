package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	sdk "github.com/dysolution/espsdk"
)

// A Bullet represents a single HTTP request that performs an operation
// against a single API endpoint. Each Bomb can contain one or multiple
// Bullets.
type Bullet struct {
	client  *sdk.Client
	Method  string `json:"method"`
	URL     string `json:"url"`
	payload sdk.Serializable
}

// Deploy sets the Bullet in motion.
func (b *Bullet) Deploy() (sdk.DeserializedObject, error) {
	switch b.Method {
	case "GET", "get":
		return b.client.Get(b.URL), nil
	case "POST", "post":
		return b.client.Create(b.URL, b.payload), nil
	case "PUT", "put":
		return b.client.Update(b.URL, b.payload), nil
	case "DELETE", "delete":
		return b.client.Delete(b.URL), nil
	}
	return sdk.DeserializedObject{}, errors.New("undefined method")
}

func (f *Bullet) String() string {
	out, err := json.MarshalIndent(f, "", "  ")
	check(err)
	return fmt.Sprintf("%s", out)
}

// A Bomb is a series of URLs and methods that represent a workflow.
type Bomb []Bullet

// Drop iterates through the Bullets within a bomb, fires all of them, and
// returns a summary of the results.
func Drop(bomb Bomb) []byte {
	var summary []byte
	for _, bullet := range bomb {
		obj, _ := bullet.Deploy()
		log.Debugf("%s", bullet)
		response, err := sdk.Marshal(obj)
		check(err)
		summary = append(summary, response...)
	}
	return summary
}

// A Raid is a collection of bombs capable of reporting summary statistics.
type Raid struct {
	StartTime time.Time `json:"start_time"`
	Payload   []Bomb    `json:"payload"`
}

// Begin iterates through the Bombs in a Raid's Payload, dropping each of
// them, and then returns a summary of the results.
func (r *Raid) Begin() []byte {
	var raidSummary []byte
	for _, bomb := range r.Payload {
		raidSummary = append(raidSummary, Drop(bomb)...)
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
func NewRaid(payload []Bomb) Raid {
	return Raid{time.Now(), payload}
}
