package airstrike

import (
	"time"

	"github.com/Sirupsen/logrus"
	sdk "github.com/dysolution/espsdk"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var log *logrus.Logger

func init() {
	log = sdk.Log
	log.Formatter = &prefixed.TextFormatter{TimestampFormat: time.RFC3339}
}

type Armed interface {
	Fire() (sdk.Result, error)
}

// An Arsenal is a collection of deployable weapons. It represents a list of
// tasks that, perfored serially, compose a workflow.
//
// Think of an arsenal as the weapons available to a single pilot within a
// squadron. Many planes can deploy their arsenal at the same time, but each
// weapon in a plane's arsenal must be deployed one at a time.
//
// For example, a common workflow would be:
//   1. list all batches
//   2. get the metadata for a batch
//   3. upload a contribution to the batch
type Arsenal struct {
	Name    string  `json:"name"`
	Weapons []Armed `json:"weapons"`
}

// Deploy sequentially fires all of the weapons within an Arsenal and reports
// the results.
func Deploy(arsenal Arsenal) ([]sdk.Result, error) {
	var results []sdk.Result
	for _, weapon := range arsenal.Weapons {
		log.Debugf("deploying %s", weapon)
		result, _ := weapon.Fire()
		results = append(results, result)
	}
	return results, nil
}
