package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	sdk "github.com/dysolution/espsdk"
)

type SimpleRaid struct {
	Bombs []SimpleBomb `json:"bombs"`
}

// A Raid is a collection of bombs capable of reporting summary statistics.
type Raid struct {
	Bombs []Bomb `json:"bombs"`
}

// Conduct concurrently drops all of the Bombs in a Raid's Payload and
// returns a collection of the results.
func (r *Raid) Conduct() ([]*sdk.Result, error) {
	logID := "Raid.Conduct"
	var allResults []*sdk.Result
	var wg sync.WaitGroup
	for i, bomb := range r.Bombs {
		wg.Add(1)
		bombID := i + 1
		go func(bombID int) ([]*sdk.Result, error) {
			defer wg.Done()

			results, err := Drop(bomb)
			if err != nil {
				log.Errorf("Raid.Conduct(): %v", err)
				return []*sdk.Result{}, err
			}

			for _, req := range results {
				log.WithFields(log.Fields{
					"bomb_id":       bombID,
					"response_time": req.VerboseResult.Duration * time.Millisecond,
					"status_code":   req.VerboseResult.Response.StatusCode,
				}).Info(logID)
				allResults = append(allResults, results...)
			}
			return allResults, nil
		}(bombID)
	}
	return allResults, nil
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
