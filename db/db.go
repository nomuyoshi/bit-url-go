package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nomuyoshi/bit-url/env"
)

const table string = "bit_url"

// DB はDynamoDBのインスタンスをもつ構造体
type DB struct {
	Instance *dynamodb.DynamoDB
}

// New はDynamoDBのインスタンスを返す
func New() DB {
	s := session.Must(session.NewSession())
	return DB{Instance: dynamodb.New(s, aws.NewConfig().WithRegion(env.Config().Region))}
}

// PutItem はDynamoDBにデータを作成する
func (db DB) PutItem(i interface{}) (*dynamodb.PutItemOutput, error) {
	// struct を*dynamodb.AttributeValueに変換してくれる
	// https://pkg.go.dev/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute?tab=doc#example-Marshal
	// https://pkg.go.dev/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute?tab=doc#MarshalMap
	item, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{Item: item}
	input.SetTableName(table)

	output, err := db.Instance.PutItem(input)

	if err != nil {
		return nil, err
	}
	return output, nil
}
