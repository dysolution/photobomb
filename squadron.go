package main

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	sdk "github.com/dysolution/espsdk"
	"github.com/dysolution/photobomb/airstrike"
)

type Squadron struct {
	wg sync.WaitGroup
}

func NewSquadron() Squadron {
	var wg sync.WaitGroup
	return Squadron{wg}
}

func (s *Squadron) bombard(ch chan sdk.Result, pilotID int, arsenal airstrike.Arsenal) {
	s.wg.Add(1)
	defer s.wg.Done()

	results, err := airstrike.Deploy(arsenal)
	if err != nil {
		log.Errorf("Raid.Conduct(): %v", err)
		ch <- sdk.Result{}
	}

	for weaponID, result := range results {
		log.WithFields(logrus.Fields{
			"pilot_id":      pilotID,
			"weapon_id":     weaponID,
			"method":        result.Verb,
			"path":          result.Path,
			"response_time": result.Duration * time.Millisecond,
			"status_code":   result.StatusCode,
		}).Info()

		ch <- result
	}
}
