all: test
	go build

vendor:: ./vendor
	go mod vendor

test: vendor
	go test ./...
