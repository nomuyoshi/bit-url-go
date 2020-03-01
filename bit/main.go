package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/speps/go-hashids"

	_ "github.com/joho/godotenv/autoload"

	"github.com/nomuyoshi/bit-url/db"
	"github.com/nomuyoshi/bit-url/env"
)

var dynamo db.DB
var config env.Env

type requestBody struct {
	url string
}

type errorResponse struct {
	message string
}

type response struct {
	url string
}

// BitURL は短縮URLともとのURLをマッピングするもの
type BitURL struct {
	Path        string `dynamodbav:"path"`
	OriginalURL string `dynamodbav:"original_url"`
}

func init() {
	dynamo = db.New()
	config = env.Config()
}

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var r requestBody

	if err := json.Unmarshal([]byte(req.Body), &r); err != nil {
		return buildResponse(
			http.StatusBadRequest,
			errorResponseBody("Failed to parse request body"),
		), err
	}

	// net/urlパッケージでパースできれば良しとする。
	pURL, err := url.Parse(r.url)
	if err != nil {
		return buildResponse(
			http.StatusBadRequest,
			errorResponseBody("Invalid url"),
		), err
	}

	hashID := generateHashID(pURL)
	bit := &BitURL{
		Path:        hashID,
		OriginalURL: r.url,
	}

	if _, err = dynamo.PutItem(bit); err != nil {
		return buildResponse(
			http.StatusInternalServerError,
			errorResponseBody("Internal Server Error"),
		), err
	}

	res := response{url: config.BaseURL + "/" + bit.Path}
	b, _ := json.Marshal(res)
	return buildResponse(
		http.StatusOK,
		string(b),
	), nil
}

func buildResponse(status int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}
}

func generateHashID(u *url.URL) string {
	hd := hashids.NewData()
	hd.Salt = config.Salt
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)
	hashID, _ := h.Encode([]int{len(u.Host), len(u.Path), int(time.Now().Unix())})

	return hashID
}

func errorResponseBody(msg string) string {
	res := errorResponse{message: msg}
	bytes, _ := json.Marshal(res)
	return string(bytes)
}
