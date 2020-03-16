package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/nomuyoshi/bit-url/db"
)

var dynamo db.DB

func init() {
	dynamo = db.New()
}

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path, err := req.PathParameters["path"]
	if err != nil {
		return response(http.StatusBadRequest, errorResponseBody(err.Error())), err
	}
	url, err := dynamo.GetItem(path)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody("Internal Server Error")
		), err
	}

	if url == "" {
		return response(http.StatusNotFound, ""), nil
	}

	return redirect(URL), nil
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}

func response(status int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}
}

func redirect(location string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers:    map[string]string{"Location": location},
	}
}