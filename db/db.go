package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nomuyoshi/bit-url/env"
)

const table string = "bit_urls"

// DB はDynamoDBのインスタンスをもつ構造体
type DB struct {
	Instance *dynamodb.DynamoDB
}

// BitURL は短縮URLともとのURLをマッピングするもの
type BitURL struct {
	Path        string `json:"path"`
	OriginalURL string `json:"original_url"`
}

// New はDynamoDBのインスタンスを返す
func New() DB {
	config := &aws.Config{Region: aws.String(env.Config().Region)}
	if env.Config().Env == "local" {
		config = config.WithEndpoint("http://dynamodb:8000")
	}
	session := session.Must(session.NewSession(config))
	return DB{Instance: dynamodb.New(session)}
}

// PutItem はDynamoDBにデータを作成する
func (db DB) PutItem(i interface{}) (*dynamodb.PutItemOutput, error) {
	// struct を*dynamodb.AttributeValueに変換してくれる
	// https://pkg.go.dev/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute?tab=doc#example-Marshal
	// https://pkg.go.dev/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute?tab=doc#MarshalMap
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}
	item, err := db.Instance.PutItem(input)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetItem はDynamoDBのデータを取得する
func (db DB) GetItem(path string) (string, error) {
	item, err := db.Instance.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"path": {
				S: aws.String(path),
			},
		},
	})

	if err != nil {
		return "", errors.Wrap(err, "failed to get bit url")
	}

	if item.Item == nil {
		return "", nil
	}

	bitURL := BitURL{}
	err = dynamodbattribute.UnmarshalMap(item.Item, &bitURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal BitURL")
	}

	return bitURL.OriginalURL, nil
}
