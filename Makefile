deps:
	go mod download

build: fmt
	cd $$GOPATH/src/github.com/caddyserver/caddy/caddy && go run build.go && cp caddy $$GOPATH/src/github.com/xuqingfeng/caddy-rate-limit/

build-docker-image:
	docker build -t xuqingfeng/caddy-rate-limit:$$(git describe --abbrev=0 --tags) . && docker build -t xuqingfeng/caddy-rate-limit:latest .

fmt:
	go fmt ./...

test: fmt
	go test -v

race-test: fmt
	go test -v -race

benchmark:
	go test -run=xxx -bench=.

benchmark-mem:
	go test -run=xxx -bench=. -benchmem

benchmark-mem-pprof:
	go test -run=xxx -bench=. -memprofile=mem.pprof

benchmark-cpu-pprof:
	go test -run=xxx -bench=. -cpuprofile=cpu.pprof
