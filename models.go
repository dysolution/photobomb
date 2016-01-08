package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sdk "github.com/dysolution/espsdk"
)

// A Flechette represents a single HTTP request that performs an operation
// against a single API endpoint. Each Bomb can contain one or multiple
// Flechettes.
type Flechette struct {
	client  *sdk.Client
	Method  string `json:"method"`
	URL     string `json:"url"`
	payload sdk.Serializable
}

// Fire makes the Flechette hit its target, hitting the endpoint with the
// described method and (optional) payload.
func (f *Flechette) Fire() (sdk.DeserializedObject, error) {
	switch f.Method {
	case "GET", "get":
		return f.client.Get(f.URL), nil
	case "POST", "post":
		return f.client.Create(f.URL, f.payload), nil
	}
	return sdk.DeserializedObject{}, errors.New("undefined method")
}

func (f *Flechette) String() string {
	out, err := json.MarshalIndent(f, "", "  ")
	check(err)
	return fmt.Sprintf("%s", out)
}

// A Bomb is a series of URLs and methods that represent a workflow.
type Bomb []Flechette

// Drop iterates through the Flechettes within a bomb, fires all of them, and
// returns a summary of the results.
func Drop(bomb Bomb) []byte {
	var summary []byte
	for _, flechette := range bomb {
		obj, _ := flechette.Fire()
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
