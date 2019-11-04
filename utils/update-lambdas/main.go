package main

import "github.com/aws/aws-sdk-go/service/cloudformation"

func LookupStack(client *cloudformation.Cloudformation, stackName string) ([]*cloudformation.StackResource, error) {
	request := cloudformation.DescribeStackResourcesInput{}
}
