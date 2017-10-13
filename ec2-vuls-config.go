package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
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
			Usage: "Filter options (default: Name=tag:vuls:scan,Values=true, Name=instance-state-name,Values=running)",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: os.Getenv("PWD") + "/config.toml",
			Usage: "Config file path",
		},
		cli.StringFlag{
			Name:  "out, o",
			Value: os.Getenv("PWD") + "/config.toml",
			Usage: "Output file path of config",
		},
		cli.BoolFlag{
			Name:  "print, p",
			Usage: "Echo the standard output instead of write into specified config file.",
		},
	}

	app.Action = func(c *cli.Context) error {

		instances, err := DescribeInstances(c.String("filters"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		config, err := LoadFile(c.String("config"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		new_config := CreateConfig(GenerateServerSection(instances), config)
		if c.Bool("print") {
			fmt.Println(string(new_config))
		} else {
			err = WriteFile(c.String("out"), new_config)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}
		return nil
	}

	app.Run(os.Args)
}
