FROM golang:1.14-alpine as builder

RUN apk add --no-cache git

COPY build/caddy.go /go/build/caddy.go
COPY build/go.mod /go/build/go.mod

RUN cd /go/build && \
    go build

FROM alpine:3.10

RUN apk add --no-cache --no-progress curl tini ca-certificates

COPY --from=builder /go/build/caddy /usr/bin/caddy

EXPOSE 2016

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["caddy"]
