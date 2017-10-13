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
	b = append(b, START+"\n"...)
	b = append(b, "# Updated "+time.Now().Format(time.RFC3339)+"\n\n"...)
	for _, instance := range instances {
		b = append(b, "[servers."+GetTagValue(instance, "Name")+"]\n"...)
		b = append(b, "host = \""+*instance.PrivateIpAddress+"\"\n"...)
		b = append(b, "\n"...)
	}
	b = append(b, END...)
	return b
}

func MergeConfig(currentConfig, newConfig []byte) []byte {

	// If it has already been created, it is rewrited.
	re := regexp.MustCompile("(?m)" + START + "[\\s\\S]*?" + END)
	if re.Match(currentConfig) {
		return re.ReplaceAll(currentConfig, newConfig)
	}

	// If it finds servers section, it is appended.
	re = regexp.MustCompile("(?m)\\[servers.*\\][\\s\\S]*")
	if re.Match(currentConfig) {
		currentConfig = append(currentConfig, newConfig...)
		return currentConfig
	}

	// In the case that it doesn't finds servers section.
	currentConfig = append(currentConfig, []byte("[servers]\n")...)
	currentConfig = append(currentConfig, newConfig...)
	return currentConfig
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
