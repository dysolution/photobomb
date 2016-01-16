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

type Squadron struct {
	wg sync.WaitGroup
}

func NewSquadron() Squadron {
	var wg sync.WaitGroup
	return Squadron{wg}
}

func (s *Squadron) bombard(ch chan sdk.Result, bombardierID int, bomb Bomb) {
	s.wg.Add(1)
	defer s.wg.Done()

	results, err := Drop(bomb)
	if err != nil {
		log.Errorf("Raid.Conduct(): %v", err)
		ch <- sdk.Result{}
	}

	for weaponID, result := range results {
		log.WithFields(logrus.Fields{
			"bombardier_id": bombardierID,
			"weapon_id":     weaponID,
			"method":        result.Verb,
			"path":          result.Path,
			"response_time": result.Duration * time.Millisecond,
			"status_code":   result.StatusCode,
		}).Infof("sitrep %v:", bombardierID)

		ch <- result
	}
}

// Conduct concurrently drops all of the Bombs in a Raid's Payload and
// returns a collection of the results.
func (r *Raid) Conduct() ([]sdk.Result, error) {
	var allResults []sdk.Result
	var reporterWg = sync.WaitGroup{}
	var ch chan sdk.Result

	for bombID, bomb := range r.Bombs {
		squadron := NewSquadron()
		go squadron.bombard(ch, bombID, bomb)
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
