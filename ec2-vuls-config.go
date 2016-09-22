package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	START = "### Generate by ec2-vuls-config ###"
	END   = "### ec2-vuls-config end ###"
)

type ServerConfig struct {
	Name string
	Host string
}

type Tag struct {
	Name   string
	Values string
}

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

		// Parse filter
		var tag Tag
		for _, set := range strings.Split(c.String("filters"), ",") {
			key := strings.Split(set, "=")
			switch key[0] {
			case "Name":
				tag.Name = key[1]
			case "Values":
				tag.Values = key[1]
			}
		}

		sess, err := session.NewSessionWithOptions(session.Options{})
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		svc := ec2.New(sess)
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String("running"),
					},
				},
				{
					Name: aws.String(tag.Name),
					Values: []*string{
						aws.String(tag.Values),
					},
				},
			},
		}
		resp, err := svc.DescribeInstances(params)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		var serverConfigs []ServerConfig
		for _, res := range resp.Reservations {
			for _, instance := range res.Instances {
				var tag_name string
				for _, t := range instance.Tags {
					if *t.Key == "Name" {
						tag_name = *t.Value
					}
				}
				// Ignore Windows instance
				if instance.Platform != nil {
					continue
				}
				serverConfigs = append(serverConfigs, ServerConfig{Name: tag_name, Host: *instance.PrivateIpAddress})
			}
		}

		// Create contents
		contents := ""
		for _, server := range serverConfigs {
			contents += "[servers." + server.Name + "]\n"
			contents += "host = \"" + server.Host + "\"\n"
			contents += "\n"
		}

		// Replace config.toml
		bs, err := ioutil.ReadFile(c.String("config"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		config := string(bs)

		re := regexp.MustCompile("(?m)" + START + "[\\s\\S]*?" + END)

		var str, generate string
		generate = START + "\n"
		generate += "# Updated " + time.Now().Format(time.RFC3339) + "\n\n"
		generate += contents
		generate += END

		if re.MatchString(config) {
			str = re.ReplaceAllString(config, generate)
		} else {
			str = config + generate
		}
		err = ioutil.WriteFile(c.String("config"), []byte(str), os.ModePerm)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}

	app.Run(os.Args)
}
