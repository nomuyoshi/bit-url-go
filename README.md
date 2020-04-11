## URL短縮サービス

```
$ curl -X POST https://mo68vfaxn5.execute-api.ap-northeast-1.amazonaws.com/Prod/bits -d '{"url":"https://news.yahoo.co.jp/topics"}'
=> {"url":"https://mo68vfaxn5.execute-api.ap-northeast-1.amazonaws.com/Prod/bits/e7ibhnzqExZ"}
```

## requirement

- [aws sam cli](https://docs.aws.amazon.com/ja_jp/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- [aws cli](https://docs.aws.amazon.com/ja_jp/streams/latest/dev/kinesis-tutorial-cli-installation.html)

## 環境構築

### build
```bash
$ git clone git@github.com:nomuyoshi/bit-url-go.git
$ cd bit-url-go
$ sam build
```

### docker
```bash
$ docker network create lambda-local
$ docker-compose up -d
```

### dynamodb-local

```bash
$ aws dynamodb create-table --cli-input-json file://local/dynamo-table.json --endpoint-url http://localhost:8000
$ aws dynamodb list-tables --endpoint-url http://localhost:8000
```

## ローカル実行

```bash
$ docker-compose up -d
$ sam local start-api --docker-network lambda-local --parameter-overrides='Salt="bit-url-local" BaseUrl="http://localhost:3000/"'
$ curl -X POST http://127.0.0.1:3000/bits -d '{"url":"https://example.com"}'
```

## ディレクトリ

- bit/ -> URL短縮処理
- redirect/ -> 短縮URLでアクセスしたときのリダイレクト処理
- db/ -> DynamoDBへのデータ保存、データ取得処理
