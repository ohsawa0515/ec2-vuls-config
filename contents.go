package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

const (
	START = "### Generate by ec2-vuls-config ###"
	END   = "### ec2-vuls-config end ###"
)

func GenerateServerSection(instances []*ec2.Instance) string {
	content := ""
	for _, instance := range instances {
		content += "[servers." + GetTagValue(instance, "Name") + "]\n"
		content += "host = \"" + *instance.PrivateIpAddress + "\"\n"
		content += "\n"
	}
	return content
}

func CreateConfig(content string, config string) string {
	re := regexp.MustCompile("(?m)" + START + "[\\s\\S]*?" + END)

	area := START + "\n"
	area += "# Updated " + time.Now().Format(time.RFC3339) + "\n\n"
	area += content
	area += END

	if re.MatchString(config) {
		return re.ReplaceAllString(config, area)
	} else {
		return config + area
	}
}

func LoadFile(path string) (string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func WriteFile(path string, content string) error {
	err := ioutil.WriteFile(path, []byte(content), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
