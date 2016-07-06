deps:
	go get -v ./...

build: format
	go build

fmt:
	go fmt ./...

test:
	go test ./...
