package buildkite

import (
	"errors"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/ec2"
	"time"
)

func EC2InstanceTags() (map[string]string, error) {
	tags := make(map[string]string)

	// Passing blank values here instructs the AWS library to look at the
	// current instances meta data for the security credentials.
	auth, err := aws.GetAuth("", "", "", time.Time{})
	if err != nil {
		return tags, errors.New(fmt.Sprintf("Error creating AWS authentication: %s", err.Error()))
	}

	// Find the current region and create a new EC2 connection
	region := aws.GetRegion(aws.InstanceRegion())
	ec2Client := ec2.New(auth, region)

	// Filter by the current machinces instance-id
	filter := ec2.NewFilter()
	filter.Add("resource-id", aws.InstanceId())

	// Describe the tags for the current instance
	resp, err := ec2Client.DescribeTags(filter)
	if err != nil {
		return tags, errors.New(fmt.Sprintf("Error downloading tags: %s", err.Error()))
	}

	// Collect the tags
	for _, tag := range resp.Tags {
		tags[tag.Key] = tag.Value
	}

	return tags, nil
}
