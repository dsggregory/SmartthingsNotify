all: test
	go build

test:
	go test ./...
