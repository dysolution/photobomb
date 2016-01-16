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
func (r *Raid) Conduct() ([]sdk.Result, error) {
	logID := "Raid.Conduct"
	var allResults []sdk.Result
	var wg sync.WaitGroup
	for bombID, bomb := range r.Bombs {
		wg.Add(1)
		go func(bombID int, bomb Bomb) ([]sdk.Result, error) {
			defer wg.Done()

			results, err := Drop(bomb)
			if err != nil {
				log.Errorf("Raid.Conduct(): %v", err)
				return []sdk.Result{}, err
			}

			for weaponID, result := range results {
				log.WithFields(logrus.Fields{
					"bomb_id":       bombID,
					"weapon_id":     weaponID,
					"method":        result.Verb,
					"path":          result.Path,
					"response_time": result.Duration * time.Millisecond,
					"status_code":   result.StatusCode,
				}).Info(logID)
				allResults = append(allResults, results...)
			}
			return allResults, nil
		}(bombID, bomb)
	}
	return allResults, nil
}

func (r *Raid) String() string {
	out, err := json.MarshalIndent(r, "", "  ")
	tableFlip(err)
	return string(out)
}
