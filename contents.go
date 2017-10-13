package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	START = "### Generate by ec2-vuls-config ###"
	END   = "### ec2-vuls-config end ###"
)

func GenerateServerSection(instances []*ec2.Instance) []byte {
	b := make([]byte, 0, 1024)
	for _, instance := range instances {
		b = append(b, "[servers."+GetTagValue(instance, "Name")+"]\n"...)
		b = append(b, "host = \""+*instance.PrivateIpAddress+"\"\n"...)
		b = append(b, "\n"...)
	}
	return b
}

func CreateConfig(content []byte, config []byte) []byte {
	re := regexp.MustCompile("(?m)" + START + "[\\s\\S]*?" + END)

	b := make([]byte, 0, 1024)
	b = append(b, START+"\n"...)
	b = append(b, "# Updated "+time.Now().Format(time.RFC3339)+"\n\n"...)
	b = append(b, content...)
	b = append(b, END...)

	// if match, return replaced contents
	if re.Match(config) {
		return re.ReplaceAll(config, b)
	}
	config = append(config, b...)
	return config
}

func LoadFile(path string) ([]byte, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func WriteFile(path string, content []byte) error {
	err := ioutil.WriteFile(path, content, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
