all: test
	go build

vendor:
	go mod vendor

test: vendor
	go test ./...
