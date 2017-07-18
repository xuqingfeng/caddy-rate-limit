deps:
	go get -v -d ./...

build: format
	go build

fmt:
	go fmt ./...

test:
	go test -v $$(go list ./... | grep -v /vendor/)

race:
	go test -v -race $$(go list ./... | grep -v /vendor/)

benchmark:
	go test -run=xxx -bench=.

benchmark-mem:
	go test -run=xxx -bench=. -benchmem

benchmark-mem-pprof:
	go test -run=xxx -bench=. -memprofile=mem.pprof

benchmark-cpu-pprof:
	go test -run=xxx -bench=. -cpuprofile=cpu.pprof
