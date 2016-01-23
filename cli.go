// Photobomb conducts workflow tests triggered by requests to its web server.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/dysolution/airstrike"
	"github.com/dysolution/espsdk"
	"github.com/dysolution/sleepwalker"
)

func main() {
	app := cli.NewApp()
	app.Name = NAME
	app.Version = VERSION
	app.Usage = "test workflows for the Getty Images ESP API"
	app.Author = "Jordan Peterson"
	app.Email = "dysolution@gmail.com"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "enable debug output",
		},
		cli.StringFlag{
			Name:   "key, k",
			Usage:  "your ESP API key",
			EnvVar: "ESP_API_KEY",
		},
		cli.StringFlag{
			Name:   "secret",
			Usage:  "your ESP API secret",
			EnvVar: "ESP_API_SECRET",
		},
		cli.StringFlag{
			Name:   "username, u",
			Usage:  "your ESP username",
			EnvVar: "ESP_USERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "your ESP password",
			EnvVar: "ESP_PASSWORD",
		},
		cli.StringFlag{
			Name:   "token, t",
			Usage:  "use an existing OAuth2 token",
			EnvVar: "ESP_TOKEN",
		},
		cli.StringFlag{
			Name:   "s3-bucket, b",
			Value:  "oregon",
			Usage:  "nearest S3 bucket = [germany|ireland|oregon|singapore|tokyo|virginia]",
			EnvVar: "S3_BUCKET",
		},
		cli.StringFlag{
			Name:   "format, f",
			Value:  "json",
			Usage:  "[json|ascii]",
			EnvVar: "PHOTOBOMB_OUTPUT_FORMAT",
		},
		cli.DurationFlag{
			Name:   "attack-interval, i",
			Value:  time.Duration(5 * time.Second),
			Usage:  "wait this long between attacks (minimum 1s)",
			EnvVar: "PHOTOBOMB_ATTACK_INTERVAL",
		},
		cli.DurationFlag{
			Name:        "warning-threshold, w",
			Value:       time.Duration(200 * time.Millisecond),
			Usage:       "log WARNINGs for long response times, e.g.: [0.2s|200ms|200000μs|200000000ns]",
			EnvVar:      "PHOTOBOMB_WARNING_THRESHOLD",
			Destination: &warningThreshold,
		},
		cli.StringFlag{
			Name:   "config, c",
			Value:  "config.json",
			Usage:  "file containing configuration (try running \"example\")",
			EnvVar: "PHOTOBOMB_CONFIG",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "suppress log output",
		},
	}
	app.Before = func(c *cli.Context) error {
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
	app.Commands = []cli.Command{
		{
			Name:  "config",
			Usage: "print a JSON representation of the config",
			Action: func(c *cli.Context) {
				out, err := json.MarshalIndent(config, "", "    ")
				tableFlip(err)
				fmt.Printf("%s\n", out)
			},
		},
		{
			Name:  "example",
			Usage: "print an example JSON configuration",
			Action: func(c *cli.Context) {
				out, err := json.MarshalIndent(ExampleConfig(), "", "    ")
				tableFlip(err)
				fmt.Printf("%s\n", out)
			},
		},
		{
			Name:  "gauge",
			Usage: "display a horizontal bar gauge of response time",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "max-width, m",
					Usage: "console width",
					Value: 80,
				},
				cli.StringFlag{
					Name:  "glyph, g",
					Usage: "character to use to build graph bars",
					Value: "=",
				},
			},
			Action: func(c *cli.Context) {
				reporter.Gauge = true
				reporter.MaxColumns = c.Int("max-width")
				reporter.Glyph = c.String("glyph")[0]
				log.Level = logrus.ErrorLevel
				serve()
			},
		},
	}
	app.Action = func(c *cli.Context) {
		serve()
	}
	app.Run(os.Args)
}
