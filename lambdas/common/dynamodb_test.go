package common

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

type MockedDynamo struct {
	dynamodbiface.DynamoDBAPI
}

func (*MockedDynamo) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if *input.TableName == "successtest" {
		out := dynamodb.PutItemOutput{}
		return &out, nil
	} else {
		return nil, errors.New("kaboom!")
	}
}

func (*MockedDynamo) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if *input.TableName == "recordstest" {
		/* prepare some test data to return */
		recordsToReturn := make([]NewReleaseEvent, 3)
		//in the "real world" Dynamo should return these records sorted into this order
		recordsToReturn[0] = NewReleaseEvent{
			Event:       "test",
			BuildId:     26,
			Branch:      "somebranch",
			DownloadUrl: "https://some/url/26",
			ProductName: "test product",
		}
		recordsToReturn[1] = NewReleaseEvent{
			Event:       "test",
			BuildId:     25,
			Branch:      "somebranch",
			DownloadUrl: "https://some/url/25",
			ProductName: "test product",
		}
		recordsToReturn[2] = NewReleaseEvent{
			Event:       "test",
			BuildId:     24,
			Branch:      "somebranch",
			DownloadUrl: "https://some/url/24",
			ProductName: "test product",
		}

		marshalledRecords := make([]map[string]*dynamodb.AttributeValue, 3)

		for ctr, entry := range recordsToReturn {
			marshalledRecords[ctr], _ = dynamodbattribute.MarshalMap(entry)
		}

		/* build the output packet */
		out := &dynamodb.QueryOutput{
			Items: marshalledRecords,
			Count: aws.Int64(3),
		}
		return out, nil

	} else if *input.TableName == "emptytest" {
		var marshalledRecords []map[string]*dynamodb.AttributeValue

		out := &dynamodb.QueryOutput{
			Items: marshalledRecords,
			Count: aws.Int64(0),
		}
		return out, nil
	} else if *input.TableName == "failtest" {
		return nil, errors.New("Kaboom!")
	} else {
		panic("query mock called without a special table name")
	}
}

func TestNewReleaseEvent_LogRelease(t *testing.T) {
	dynamoClient := &MockedDynamo{}

	evt := NewReleaseEvent{
		Event:       "test",
		BuildId:     1234,
		Branch:      "master",
		DownloadUrl: "https://some/url",
		ProductName: "test product",
	}

	err := evt.LogRelease(dynamoClient, "successtest")
	if err != nil {
		t.Errorf("put test should have succeeded but got %s", err)
	}

	shouldErr := evt.LogRelease(dynamoClient, "failtest")
	if shouldErr == nil {
		t.Errorf("put test should have failed but returned no error")
	}
}

func TestMostRecentRelease(t *testing.T) {
	dynamoClient := &MockedDynamo{}

	result, err := MostRecentRelease(dynamoClient, "recordstest", "test product", "somebranch")
	if err != nil {
		t.Errorf("read test for records should have succeeded but got %s", err)
	}
	if result == nil {
		t.Errorf("read test for records should have returned a record but returned nil")
	}
	if result.ProductName != "test product" {
		t.Errorf("returned product name was wrong")
	}
	if result.DownloadUrl != "https://some/url/26" {
		t.Errorf("returned download url was wrong, got %s", result.DownloadUrl)
	}
	if result.Branch != "somebranch" {
		t.Errorf("returned branch was wrong, got %s", result.Branch)
	}
	if result.BuildId != 26 {
		t.Errorf("returned build id was wrong, got %d", result.BuildId)
	}

	emptyResult, emptyErr := MostRecentRelease(dynamoClient, "emptytest", "test product", "somebranch")
	if emptyErr != nil {
		t.Errorf("empty result test should have succeeded but got %s", err)
	}
	if emptyResult != nil {
		t.Errorf("empty result test should have returned nil but got %s", spew.Sprint(*emptyResult))
	}

	failedResult, failedErr := MostRecentRelease(dynamoClient, "failtest", "test product", "somebranch")
	if failedErr == nil {
		t.Errorf("failure test should have failed but got nil error")
	}
	if failedResult != nil {
		t.Errorf("failure test should have returned nil but got %s", spew.Sprint(*failedResult))
	}
}
