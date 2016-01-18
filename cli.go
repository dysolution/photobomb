// Photobomb conducts workflow tests triggered by requests to its web server.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
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
		cli.StringFlag{
			Name:   "config, c",
			Value:  "config.json",
			Usage:  "file containing configuration (try running \"example\")",
			EnvVar: "PHOTOBOMB_CONFIG",
		},
	}
	app.Before = func(c *cli.Context) error {
		client = sleepwalker.GetClient(
			c.String("key"),
			c.String("secret"),
			c.String("username"),
			c.String("password"),
			espsdk.OAuthEndpoint,
			espsdk.ESPAPIRoot,
			log,
		)
		log.Debugf("client, created from environment: %v", client)

		if c.Bool("debug") == true {
			log.Level = logrus.DebugLevel
		} else {
			log.Level = logrus.InfoLevel
		}

		token = sleepwalker.Token(c.String("token"))

		if strings.ToLower(c.String("format")) == "json" {
			log.Formatter = &logrus.JSONFormatter{}
		}

		config = loadConfig(c.String("config"))
		configJSON, err := json.Marshal(config)
		tableFlip(err)
		log.Debugf("configuration: %s", configJSON)

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
	}
	app.Action = func(c *cli.Context) { httpd() }

	app.Run(os.Args)
}
