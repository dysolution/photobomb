package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/dysolution/airstrike"
	"github.com/dysolution/espsdk"
	"github.com/dysolution/sleepwalker"
)

func appBefore(c *cli.Context) error {
	desc := "cli.app.Before"

	switch {
	case c.Bool("debug"):
		log.Level = logrus.DebugLevel
	case c.Bool("quiet"):
		log.Level = logrus.ErrorLevel
	default:
		log.Level = logrus.InfoLevel
	}

	client = sleepwalker.GetClient(
		c.String("key"),
		c.String("secret"),
		c.String("username"),
		c.String("password"),
		espsdk.OAuthEndpoint,
		espsdk.ESPAPIRoot,
		log,
	)

	// set up the reporter for logging and console output
	reporter = airstrike.Reporter{
		CountGoroutines:  false, // caution: uses package runtime
		Logger:           log,
		URLInvariant:     espsdk.APIInvariant,
		WarningThreshold: warningThreshold,
	}

	token = sleepwalker.Token(c.String("token"))

	if strings.ToLower(c.String("format")) == "json" {
		log.Formatter = &logrus.JSONFormatter{}
	}

	cliInterval := int(c.Duration("attack-interval") / time.Duration(1*time.Second))
	if cliInterval != 0 {
		interval = cliInterval
	}

	config = loadConfig(c.String("config"))
	cfgJSON, err := json.Marshal(config)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": "unable to marshal config",
		}).Error(desc)
	}
	log.WithFields(logrus.Fields{"config": string(cfgJSON)}).Debug(desc)

	return nil
}
