package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"log"
)

/**
logs a new release to the database
arguments:
 - client: an instance of Dynamodb client or a mock
 - tableName: table name to write to. Client must have PutObject permission for this
*/
func (ev *NewReleaseEvent) LogRelease(client dynamodbiface.DynamoDBAPI, tableName string) error {
	attributeValues, marshalErr := dynamodbattribute.MarshalMap(ev)
	if marshalErr != nil {
		log.Printf("Could not marshal data into dynamo format: %s\n", marshalErr)
		return marshalErr
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(tableName),
	}

	_, putErr := client.PutItem(input)
	if putErr != nil {
		log.Printf("Could not put item to Dynamo table %s: %s\n", tableName, putErr)
		return putErr
	}

	log.Printf("Successfully wrote data to dynamo %s\n", tableName)
	return nil
}

/**
retrieve a record for the most recent release of productName
arguments:
    - client: an instance of Dynamodb client or a mock
    - tableName: a string of the dynamo table name
    - productName: product name to filter on
    - branch: branch to filter on
returns either:
	- nil and an error if an error occurred
    - a pointer to a NewReleaseEvent record and nil if a record was found
    - nil and nil if no record was found and there was not an error
*/
func MostRecentRelease(client dynamodbiface.DynamoDBAPI, tableName string, productName string, branch string) (*NewReleaseEvent, error) {
	scanForward := false //we want to start with the highest number

	parameters := map[string]*dynamodb.AttributeValue{
		":nameSubst":   {S: aws.String(productName)},
		":branchSubst": {S: aws.String(branch)},
	}

	qInput := &dynamodb.QueryInput{
		TableName:                 &tableName,
		KeyConditionExpression:    aws.String("productName=:nameSubst"),
		ExpressionAttributeValues: parameters,
		ScanIndexForward:          &scanForward,
		FilterExpression:          aws.String("branch=:branchSubst"),
	}

	results, scanErr := client.Query(qInput)
	if scanErr != nil {
		log.Printf("Could not perform table query: %s", scanErr)
		return nil, scanErr
	}

	if *results.Count > 0 {
		objectsList := make([]NewReleaseEvent, *results.Count)

		unmarshalErr := dynamodbattribute.UnmarshalListOfMaps(results.Items, &objectsList)
		if unmarshalErr != nil {
			log.Printf("Could not unmarshal data from database: %s", unmarshalErr)
			return nil, unmarshalErr
		}
		return &objectsList[0], nil
	} else {
		log.Printf("No results found for productName %s", productName)
		return nil, nil
	}
}
