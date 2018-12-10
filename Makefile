test:
	go test ./...
	
all:: test
	go build
