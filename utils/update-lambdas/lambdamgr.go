package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
)

type LambdaMgr struct {
	ResourcesList       []*cloudformation.StackResource
	client              *cloudformation.CloudFormation
	lambaClient         *lambda.Lambda
	stackInfo           *cloudformation.Stack
	LambdaCodeLocations map[string]string
}

/**
return a new LambdaMgr instance, this will allocate the client too
*/
func NewLambdaMgr() *LambdaMgr {
	newMgr := LambdaMgr{}

	//set up an AWS session to communicate with Dynamo
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	newMgr.client = cloudformation.New(sess)
	newMgr.lambaClient = lambda.New(sess)
	return &newMgr
}

/**
look up the given lambda function and return its source path
*/
func (self *LambdaMgr) LookupLambdaSourcePath(lambdaName *string) (*string, error) {
	req := lambda.GetFunctionInput{
		FunctionName: lambdaName,
	}

	result, err := self.lambaClient.GetFunction(&req)
	if err != nil {
		log.Printf("Could not look up lambda %s: %s\n", lambdaName, err)
		return nil, err
	}

	return result.Code.Location, nil
}

/**
retrieves the stack resources associated with the given stack
*/
func (self *LambdaMgr) LookupStack(stackName *string) error {
	mdRequest := cloudformation.DescribeStacksInput{
		StackName: stackName,
	}

	mdOut, mdErr := self.client.DescribeStacks(&mdRequest)
	if mdErr != nil {
		return mdErr
	}

	self.stackInfo = mdOut.Stacks[0]
	request := cloudformation.DescribeStackResourcesInput{
		StackName: stackName,
	}

	outpt, err := self.client.DescribeStackResources(&request)
	if err != nil {
		return err
	}
	self.ResourcesList = outpt.StackResources

	self.LambdaCodeLocations = make(map[string]string)
	for _, entry := range self.ResourcesList {
		if *entry.ResourceType == "AWS::Lambda::Function" {
			resourceId := entry.PhysicalResourceId
			codeLoc, lookupErr := self.LookupLambdaSourcePath(resourceId)
			if lookupErr != nil {
				return err
			}
			self.LambdaCodeLocations[*resourceId] = *codeLoc
		}
	}
	return nil
}

/**
returns a list of the PhysicalResourceId of each lambda function
*/
func (self *LambdaMgr) LambdaFunctions() []*string {
	temp := make([]*string, len(self.ResourcesList))
	i := 0
	for _, entry := range self.ResourcesList {
		if *entry.ResourceType == "AWS::Lambda::Function" {
			temp[i] = entry.PhysicalResourceId
			i += 1
		}
	}
	return temp[:i]
}

func (self *LambdaMgr) Description() *string {
	if self.stackInfo == nil {
		return nil
	}
	return self.stackInfo.Description
}

func (self *LambdaMgr) DeployablesBucket() *string {
	if self.stackInfo == nil {
		return nil
	}

	for _, param := range self.stackInfo.Parameters {
		if *param.ParameterKey == "DeployablesBucket" {
			return param.ParameterValue
		}
	}
	return nil
}

func (self *LambdaMgr) AppStackStage() (string, string, string) {
	if self.stackInfo == nil {
		return "", "", ""
	}

	var app string
	var stack string
	var stage string

	for _, param := range self.stackInfo.Parameters {
		if *param.ParameterKey == "App" {
			app = *param.ParameterValue
		} else if *param.ParameterKey == "Stack" {
			stack = *param.ParameterValue
		} else if *param.ParameterKey == "Stage" {
			stage = *param.ParameterValue
		}
	}
	return app, stack, stage
}

func (self *LambdaMgr) UpdateLambdaCode(funcName *string, uploadBucket *string, newCodePath *string) error {
	req := lambda.UpdateFunctionCodeInput{
		FunctionName: funcName,
		Publish:      aws.Bool(true),
		S3Bucket:     uploadBucket,
		S3Key:        newCodePath,
	}

	_, err := self.lambaClient.UpdateFunctionCode(&req)
	return err
}
