.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./build

build:
	GOOS=linux GOARCH=amd64 go build -o ./bit-url/bit ./bit
