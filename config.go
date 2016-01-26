package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/dysolution/airstrike"
	"github.com/dysolution/espsdk"
	"github.com/dysolution/sleepwalker"
	"github.com/spf13/viper"
)

type Config struct {
	Output struct {
		Format string `json:"format"`
		Quiet  bool   `json:"quiet"`
		Level  string `json:"level"`
	} `json:"output"`
	Mission airstrike.Mission `json:"mission"`
}

func getConfig() Config {

	viper.SetEnvPrefix(NAME)
	viper.AutomaticEnv()

	var cfg Config

	viper.SetConfigType("json")
	viper.SetConfigName("viperconfig")
	viper.AddConfigPath("/etc/" + NAME + "/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.Set("mission.inception", time.Now())
	viper.SetDefault("mission.interval", 1)
	viper.SetDefault("output.format", "json")

	viper.Unmarshal(&cfg)
	out, _ := json.MarshalIndent(cfg, "", "    ")
	fmt.Printf("%s", out)
	return cfg
}

func appBefore(c *cli.Context) error {
	desc := "cli.app.Before"

	cfg = getConfig()

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

	if viper.GetString("format") == "json" {
		log.Formatter = &logrus.JSONFormatter{}
	}

	cliInterval := float64(c.Duration("attack-interval") / time.Duration(time.Millisecond))
	if cliInterval != 0 {
		cfg.Mission.Interval = cliInterval
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
