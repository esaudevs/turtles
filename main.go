package main

import (
	"context"
	"os"
	"strings"

	"github.com/esaudevs/turtles/awsgo"
	"github.com/esaudevs/turtles/bd"
	"github.com/esaudevs/turtles/handlers"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(RunLambda)
}

func RunLambda(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	awsgo.InitAWS()

	if !ParamsAreValid() {
		panic("Error getting params, SecretName and UrlPrefix are needed.")
	}

	var res *events.APIGatewayProxyResponse
	prefix := os.Getenv("UrlPrefix")
	path := strings.Replace(request.RawPath, prefix, "", -1)
	method := request.RequestContext.HTTP.Method
	body := request.Body
	header := request.Headers

	bd.ReadSecret()

	status, message := handlers.Handlers(path, method, body, header, request)

	headersResp := map[string]string{
		"Content-Type": "application/json",
	}

	res = &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(message),
		Headers:    headersResp,
	}

	return res, nil
}

func ParamsAreValid() bool {
	_, couldGetParams := os.LookupEnv("SecretName")
	if !couldGetParams {
		return couldGetParams
	}

	_, couldGetParams = os.LookupEnv("UrlPrefix")
	if !couldGetParams {
		return couldGetParams
	}

	return couldGetParams
}
