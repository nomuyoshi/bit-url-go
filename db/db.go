package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	_ "github.com/joho/godotenv/autoload"
)

const (
	region string = os.Getenv("REGION")
	table  string = "bit_url"
)

// DB はDynamoDBのインスタンスをもつ構造体
type DB struct {
	Instance *dynamodb.DynamoDB
}

// BitUrl は短縮URLともとのURLをマッピングするもの
type BitUrl struct {
	Path        string
	OriginalUrl string
}

// New はDynamoDBのインスタンスを返す
func New() DB {
	session = session.Must(session.NewSession())
	return DB{Instance: dynamodb.New(mySession, aws.NewConfig().WithRegion(region))}
}

// PutItem はDynamoDBにデータを作成する
func (db DB) PutItem(b *BitUrl) (*dynamodb.PutItemOutput, error) {
	// struct を*dynamodb.AttributeValueに変換してくれる
	// https://pkg.go.dev/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute?tab=doc#example-Marshal
	item, err := dynamodbattribute.Marshal(b)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: table,
	}

	output, err := db.Instance.PutItem(input)

	if err != nil {
		return nil, err
	}
	return output, nil
}
