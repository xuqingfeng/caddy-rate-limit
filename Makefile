deps:
	go get -v -d ./...

build: format
	go build

fmt:
	go fmt ./...

test:
	go test -v $$(go list ./... | grep -v /vendor/)

benchmark:
	go test -run=xxx -bench=.
