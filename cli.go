// Photobomb conducts workflow tests triggered by requests to its web server.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
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
			Name:  "format, f",
			Value: "json",
			Usage: "[json|ascii]",
			// EnvVar: "PHOTOBOMB_OUTPUT_FORMAT",
		},
		cli.DurationFlag{
			Name:   "attack-interval, i",
			Value:  time.Duration(5000 * time.Millisecond),
			Usage:  "wait n ms between attacks",
			EnvVar: "PHOTOBOMB_INTERVAL",
		},
		cli.DurationFlag{
			Name:        "warning-threshold, w",
			Value:       time.Duration(200 * time.Millisecond),
			Usage:       "log WARNINGs for long response times, e.g.: [0.2s|200ms|200000μs|200000000ns]",
			EnvVar:      "PHOTOBOMB_WARNING_THRESHOLD",
			Destination: &warningThreshold,
		},
		cli.BoolFlag{
			Name:   "quiet, q",
			Usage:  "suppress log output",
			EnvVar: "PHOTOBOMB_QUIET",
		},
	}
	app.Before = appBefore
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
