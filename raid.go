package main

import (
	"encoding/json"
	"sync"

	sdk "github.com/dysolution/espsdk"
	"github.com/dysolution/photobomb/airstrike"
)

type SimpleRaid struct {
	Arsenals []struct {
		Name    string `json:"name"`
		Weapons []struct {
			Name string `json:"name"`
		} `json:"weapons"`
	} `json:"planes"`
}

// A Raid is a collection of bombs capable of reporting summary statistics.
type Raid struct {
	Arsenals []airstrike.Arsenal `json:"planes"`
}

// Conduct concurrently drops all of the Bombs in a Raid's Payload and
// returns a collection of the results.
func (r *Raid) Conduct() ([]sdk.Result, error) {
	var allResults []sdk.Result
	var reporterWg = sync.WaitGroup{}
	var ch chan sdk.Result

	for arsenalID, arsenal := range r.Arsenals {
		squadron := airstrike.NewSquadron()
		go squadron.Bombard(ch, arsenalID, arsenal)
		go func() {
			reporterWg.Add(1)
			result := <-ch
			allResults = append(allResults, result)
		}()
	}
	return allResults, nil
}

func (r *Raid) String() string {
	out, err := json.MarshalIndent(r, "", "  ")
	tableFlip(err)
	return string(out)
}

// NewRaid initializes and returns a Raid, . It should be used in lieu of Raid literals.
func NewRaid(arsenals ...airstrike.Arsenal) Raid {
	var payload []airstrike.Arsenal
	for _, arsenal := range arsenals {
		payload = append(payload, arsenal)
	}
	return Raid{Arsenals: payload}
}
