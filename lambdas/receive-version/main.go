package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/fredex42/downloadmanager-versions-api/lambdas/common"
	"log"
	"os"
	"time"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing request with ID %s.\n", request.RequestContext.RequestID)
	log.Printf("Body size is %d\n", len(request.Body))

	tableName := os.Getenv("DYNAMO_TABLE_NAME")

	var releaseEvent common.NewReleaseEvent

	log.Printf("Got request body '%s'", request.Body)
	unmarshalErr := json.Unmarshal([]byte(request.Body), &releaseEvent)
	if unmarshalErr != nil {
		log.Printf("Could not unmarshal request body: %s\n", unmarshalErr)
		return events.APIGatewayProxyResponse{StatusCode: 400}, unmarshalErr
	}

	validationErr := releaseEvent.Validate()
	if validationErr != nil {
		log.Printf("Incoming JSON was not valid: %s\n", validationErr)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: validationErr.Error()}, nil
	}

	releaseEvent.Timestamp = time.Now().Format(time.RFC3339)

	//set up an AWS session to communicate with Dynamo
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := dynamodb.New(sess)

	putErr := releaseEvent.LogRelease(client, tableName)
	if putErr != nil {
		return events.APIGatewayProxyResponse{Body: "Could not communicate with database", StatusCode: 500}, errors.New("Could not write record to Dynamo: " + putErr.Error())
	} else {
		return events.APIGatewayProxyResponse{StatusCode: 201}, nil
	}

}

func main() {
	lambda.Start(HandleRequest)
}
