deps:
	go get -v -d ./...

build: format
	go build

fmt:
	go fmt ./...

test:
	go test $(go list ./... | grep -v /vendor/)
