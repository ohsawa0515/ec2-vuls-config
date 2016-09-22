package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ec2-vuls-config"
	app.Usage = "Generate Vuls config by filtering the Amazon EC2 information."
	app.Author = "Shuichi Ohsawa"
	app.Email = "ohsawa0515@gmail.com"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "filters, f",
			Value: "Name=tag:Vuls-Scan,Values=True",
			Usage: "Filtering ec2 tag",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: os.Getenv("PWD") + "/config.toml",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Action = func(c *cli.Context) error {

		tag := ParseFilter(c.String("filters"))
		instances, err := DescribeInstances(tag)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		config, err := LoadFile(c.String("config"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		new_config := CreateConfig(GenerateServerSection(instances), config)
		err = WriteFile(c.String("config"), new_config)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}

	app.Run(os.Args)
}
