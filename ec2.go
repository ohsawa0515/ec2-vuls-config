package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

type Tag struct {
	Name   string
	Values string
}

func ParseFilter(filters string) Tag {
	var tag Tag
	for _, set := range strings.Split(filters, ",") {
		key := strings.Split(set, "=")
		switch key[0] {
		case "Name":
			tag.Name = key[1]
		case "Values":
			tag.Values = key[1]
		}
	}
	return tag
}

func generateSession() (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{})
}

func DescribeInstances(tag Tag) ([]*ec2.Instance, error) {

	sess, err := generateSession()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	if len(resp.Reservations) == 0 {
		return []*ec2.Instance{}, nil
	}
	instances := make([]*ec2.Instance, 0)
	for _, res := range resp.Reservations {
		for _, instance := range res.Instances {

			// Ignore Windows instance
			if instance.Platform != nil {
				continue
			}

			instances = append(instances, instance)
		}
	}
	return instances, nil
}

func GetTagValue(instance *ec2.Instance, tag_name string) string {
	for _, t := range instance.Tags {
		if *t.Key == tag_name {
			return *t.Value
		}
	}
	return ""
}
