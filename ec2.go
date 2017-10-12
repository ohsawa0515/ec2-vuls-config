package main

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// parseFilter parse filter command line option.
func parseFilter(filters string) []*ec2.Filter {

	// filters e.g. "Name=tag:Foo,Values=Bar Name=instance-type,Values=m1.small"
	var ec2Filters []*ec2.Filter

	re := regexp.MustCompile(`Name=(.+),Values=(.+)`)
	for _, i := range strings.Fields(filters) {
		matches := re.FindAllStringSubmatch(i, -1)
		ec2Filters = append(ec2Filters, &ec2.Filter{
			Name: aws.String(matches[0][1]),
			Values: []*string{
				aws.String(matches[0][2]),
			},
		})
	}
	return ec2Filters
}

// generateFilter generates filter for list EC2 instances.
func generateFilter(filters string) []*ec2.Filter {
	var ec2Filters []*ec2.Filter
	if len(filters) > 0 {
		ec2Filters = parseFilter(filters)
	}
	// Default scan tag (vuls:scan)
	ec2Filters = append(ec2Filters, &ec2.Filter{
		Name: aws.String("tag:vuls:scan"),
		Values: []*string{
			aws.String("true"),
		},
	})
	// Only running instances
	ec2Filters = append(ec2Filters, &ec2.Filter{
		Name: aws.String("instance-state-name"),
		Values: []*string{
			aws.String("running"),
		},
	})
	return ec2Filters
}

// generateSession generate session.
func generateSession() (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{})
}

// DescribeInstances return list of EC2 instances.
func DescribeInstances(filters string) ([]*ec2.Instance, error) {

	sess, err := generateSession()
	if err != nil {
		return nil, err
	}
	svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		Filters: generateFilter(filters),
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

// GetTagValue returns value of EC2 tag.
func GetTagValue(instance *ec2.Instance, tag_name string) string {
	for _, t := range instance.Tags {
		if *t.Key == tag_name {
			return *t.Value
		}
	}
	return ""
}
