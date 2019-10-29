package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

/**
logs a new release to the database
*/
func (ev *NewReleaseEvent) LogRelease(tableName string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := dynamodb.New(sess)

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
