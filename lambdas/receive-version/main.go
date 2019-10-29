package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fredex42/downloadmanager-versions-api/lambdas/common"
	"log"
	"os"
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
		return events.APIGatewayProxyResponse{StatusCode: 400}, errors.New("Could not unmarshal request body")
	}

	putErr := releaseEvent.LogRelease(tableName)
	if putErr != nil {
		return events.APIGatewayProxyResponse{Body: putErr.Error(), StatusCode: 500}, errors.New("Not yet implemented")
	} else {
		return events.APIGatewayProxyResponse{StatusCode: 201}, nil
	}

}

func main() {
	lambda.Start(HandleRequest)
}
