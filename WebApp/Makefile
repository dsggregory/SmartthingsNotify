all: test
	go build

vendor:
	go mod vendor

test: vendor
	TZ=UTC go test ./...
