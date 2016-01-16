package main

import (
	sdk "github.com/dysolution/espsdk"
)

type SimpleBomb struct {
	Name    string         `json:"name"`
	Bullets []SimpleBullet `json:"bullets"`
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
// returns a slice of the results.
func Drop(bomb Bomb) ([]*sdk.Result, error) {
	var results []*sdk.Result
	for _, bullet := range bomb.Bullets {
		result, err := bullet.Deploy()
		if err != nil {
			log.Errorf("Drop: %v", err)
			return []*sdk.Result{}, err
		}
		results = append(results, result)
	}
	return results, nil
}
