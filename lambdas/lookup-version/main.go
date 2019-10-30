package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/fredex42/downloadmanager-versions-api/lambdas/common"
	"log"
	"os"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var searchReq common.SearchRequest

	tableName := os.Getenv("DYNAMO_TABLE_NAME")

	unmarshalErr := json.Unmarshal([]byte(request.Body), &searchReq)
	if unmarshalErr != nil {
		return events.APIGatewayProxyResponse{Body: "Could not understand request body", StatusCode: 400}, nil
	}

	//set up an AWS session to communicate with Dynamo
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := dynamodb.New(sess)

	var outArrayLen int
	if searchReq.AlwaysShowMaster == true {
		outArrayLen = 2
	} else {
		outArrayLen = 1
	}

	results := make([]*common.NewReleaseEvent, outArrayLen)

	branchRecord, getErr := common.MostRecentRelease(client, tableName, searchReq.ProductName, searchReq.Branch)
	if getErr != nil {
		log.Printf("Could not get data from database: %s", getErr)
		return events.APIGatewayProxyResponse{Body: "Could not get info from database", StatusCode: 500}, nil
	}
	results[0] = branchRecord

	if searchReq.AlwaysShowMaster == true {
		masterRecord, getErr := common.MostRecentRelease(client, tableName, searchReq.ProductName, "master")
		if getErr != nil {
			log.Printf("Could not get data from database: %s", getErr)
			return events.APIGatewayProxyResponse{Body: "Could not get info from database", StatusCode: 500}, nil
		}
		results[1] = masterRecord
	}

	if results[0] == nil {
		return events.APIGatewayProxyResponse{Body: "Nothing found for product and branch", StatusCode: 404}, nil
	}

	output, marshalErr := json.Marshal(results)

	if marshalErr != nil {
		log.Printf("Could not marshal final results: %s", getErr)
		return events.APIGatewayProxyResponse{Body: "Could not marshal final response", StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(output), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
