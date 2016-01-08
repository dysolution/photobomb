package main

import (
	"errors"
	_ "fmt"
	"time"

	sdk "github.com/dysolution/espsdk"
)

// A Flechette represents a single HTTP request that performs an operation
// against a single API endpoint. Each Bomb can contain one or multiple
// Flechettes.
type Flechette struct {
	Client  sdk.Client
	Method  string `json:"method"`
	URL     string `json:"url"`
	Payload sdk.Serializable
}

// Fire makes the Flechette hit its target, hitting the endpoint with the
// described method and (optional) payload.
func (o *Flechette) Fire() (sdk.DeserializedObject, error) {
	switch o.Method {
	case "GET", "get":
		return o.Client.Get(o.URL), nil
	case "POST", "post":
		return o.Client.Create(o.URL, o.Payload), nil
	}
	return sdk.DeserializedObject{}, errors.New("undefined method")
}

// A Bomb is a series of URLs and methods that represent a workflow.
type Bomb []Flechette

func Drop(b Bomb) []byte {
	var summary []byte
	for _, f := range b {
		obj, _ := f.Fire()
		response, err := sdk.Marshal(obj)
		check(err)
		summary = append(summary, response...)
	}
	return summary
}

type Raid struct {
	StartTime time.Time
	Payload   []Bomb
}

func (r *Raid) Begin() []byte {
	var raidSummary []byte
	for _, bomb := range r.Payload {
		raidSummary = append(raidSummary, Drop(bomb)...)
	}
	return raidSummary
}

func (r *Raid) Duration() time.Duration {
	return time.Now().Sub(r.StartTime)
}

func NewRaid(payload []Bomb) Raid {
	return Raid{time.Now(), payload}
}
