package main

import (
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/dysolution/espsdk"
)

type SimpleBullet struct {
	Name string `json:"name"`
}

// A Bullet represents a single HTTP request that performs an operation
// against a single API endpoint. Each Bomb can contain one or multiple
// Bullets, which are deployed serially.
type Bullet struct {
	client  *sdk.Client
	Name    string         `json:"name"`
	Method  string         `json:"method"`
	URL     string         `json:"url"`
	Payload sdk.RESTObject `json:"payload,omitempty"`
}

func (b *Bullet) handler(fn func(sdk.Findable) (sdk.Result, error)) (sdk.Result, error) {
	result, err := fn(b.Payload)
	if err != nil {
		log.Errorf("%s.Deploy %s: %v", b.Name, b.Method, err)
		return sdk.Result{}, err
	}
	result.Log().Debugf("%s.Deploy", b.Name)
	return result, nil
}

// Deploy sets the Bullet in motion.
func (b Bullet) Deploy() (sdk.Result, error) {
	switch b.Method {
	case "GET", "get":
		return b.handler(b.client.VerboseGet)
	case "POST", "post":
		return b.handler(b.client.VerboseCreate)
	case "PUT", "put":
		return b.handler(b.client.VerboseUpdate)
	case "DELETE", "delete":
		return b.handler(b.client.VerboseDelete)
	}
	msg := fmt.Sprintf("%s.Deploy: undefined method: %s", b.Name, b.Method)
	return sdk.Result{}, errors.New(msg)
}

func (b *Bullet) String() string {
	out, err := json.MarshalIndent(b, "", "  ")
	check(err)
	return fmt.Sprintf("%s", out)
}
