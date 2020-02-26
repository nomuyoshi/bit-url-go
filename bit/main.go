package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/speps/go-hashids"

	_ "github.com/joho/godotenv/autoload"
)

type requestBody struct {
	url string
}

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var r requestBody

	if err := json.Unmarshal([]byte(req.Body), &r); err != nil {
		return nil, fmt.Errorf("[Error]: failed to parse request body. %s", err)
	}

	// net/urlパッケージでパースできれば良しとする。
	pURL, err := url.Parse(r.url)
	if err != nil {
		return nil, fmt.Errorf("Error: invalid url. %s", err)
	}

	hashID = generateHashID(pURL)

}

func generateHashID(u *url.URL) string {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("HASHID_SALT")
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)
	hashID, _ := h.Encode([]int{len(pUrl.Host), len(pURL.Path), int(time.Now().Unix())})

	return hashID
}
