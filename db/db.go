package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

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
	session := session.Must(session.NewSession(&aws.Config{Region: aws.String(env.Config().Region)}))
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
