package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
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

func (b *Bullet) handler(fn func(sdk.RESTObject) (*sdk.FulfilledRequest, error)) (*sdk.FulfilledRequest, error) {
	req, err := fn(b.Payload)
	if err != nil {
		log.Errorf("%s.Deploy %s: %v", b.Name, b.Method, err)
		return &sdk.FulfilledRequest{}, err
	}
	log.WithFields(req.Stats()).Infof("%s.Deploy", b.Name)
	return req, nil
}

// Deploy sets the Bullet in motion.
func (b *Bullet) Deploy() (*sdk.FulfilledRequest, error) {
	switch b.Method {
	case "GET", "get":
		return b.handler(b.client.VerboseGet)
	case "POST", "post":
		return b.handler(b.client.VerboseCreate)
	case "PUT", "put":
		return b.handler(b.client.VerboseUpdate)
	case "DELETE", "delete":
		b.client.DeleteFromObject(b.Payload)
		return &sdk.FulfilledRequest{}, nil
	}
	msg := fmt.Sprintf("%s.Deploy: undefined method: %s", b.Name, b.Method)
	return &sdk.FulfilledRequest{}, errors.New(msg)
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
// returns a slice of the results.
func Drop(bomb Bomb) ([]*sdk.FulfilledRequest, error) {
	var results []*sdk.FulfilledRequest
	for _, bullet := range bomb.Bullets {
		result, err := bullet.Deploy()
		if err != nil {
			log.Errorf("Drop: %v", err)
			return []*sdk.FulfilledRequest{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

// A Raid is a collection of bombs capable of reporting summary statistics.
type Raid struct {
	Bombs []Bomb `json:"bombs"`
}

// Conduct concurrently drops all of the Bombs in a Raid's Payload and
// returns a collection of the results.
func (r *Raid) Conduct() ([]*sdk.FulfilledRequest, error) {
	logID := "Raid.Conduct"
	var allResults []*sdk.FulfilledRequest
	var wg sync.WaitGroup
	for i, bomb := range r.Bombs {
		wg.Add(1)
		bombID := i + 1
		go func(bombID int) ([]*sdk.FulfilledRequest, error) {
			defer wg.Done()

			results, err := Drop(bomb)
			if err != nil {
				log.Errorf("Raid.Conduct(): %v", err)
				return []*sdk.FulfilledRequest{}, err
			}

			for _, req := range results {
				log.WithFields(log.Fields{
					"bomb_id":       bombID,
					"response_time": req.Result.Duration * time.Millisecond,
					"status_code":   req.Result.Response.StatusCode,
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
