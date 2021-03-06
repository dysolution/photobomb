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
	Mission *airstrike.Mission `json:"mission"`
}

func getConfig() (cfg Config) {
	viper.SetEnvPrefix(NAME)
	viper.AutomaticEnv()

	viper.SetConfigType("json")
	viper.SetConfigName("viperconfig")
	viper.AddConfigPath("/etc/" + NAME + "/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.Set("mission.inception", time.Now())
	viper.SetDefault("output.format", "json")

	viper.Unmarshal(&cfg)
	return cfg
}

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

	getConfig()
	cfg.Mission = airstrike.NewMission(log)

	client = sleepwalker.GetClient(&sleepwalker.Config{
		Credentials: &sleepwalker.Credentials{
			APIKey:    c.String("key"),
			APISecret: c.String("secret"),
			Username:  c.String("username"),
			Password:  c.String("password"),
		},
		OAuthEndpoint: espsdk.OAuthEndpoint,
		APIRoot:       espsdk.SandboxAPI,
		Logger:        log,
	})

	cfg.Mission.Enabled = true

	cliInterval := float64(c.Duration("attack-interval") / time.Duration(time.Millisecond))
	if cliInterval > 0 {
		cfg.Mission.Interval = cliInterval
	}

	if c.Duration("warning-threshold") == 0 {
		warningThreshold = time.Duration(cfg.Mission.Interval) * time.Millisecond
	}

	// set up the reporter for logging and console output
	cfg.Mission.Reporter.URLInvariant = espsdk.APIInvariant
	cfg.Mission.Reporter.WarningThreshold = warningThreshold

	token = sleepwalker.Token(c.String("token"))

	if viper.GetString("format") == "json" {
		log.Formatter = &logrus.JSONFormatter{}
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
