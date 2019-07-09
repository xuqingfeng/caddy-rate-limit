FROM golang:1.12.7-alpine as builder

ARG CADDY_VERSION="1.0.1"

ENV GO111MODULE=on

RUN apk add --no-cache git

COPY caddy.go /go/build/caddy.go
COPY go.mod /go/build/go.mod

RUN cd /go/build && \
    go build

FROM alpine:3.10

RUN apk add --no-cache --no-progress curl tini ca-certificates

COPY --from=builder /go/build/caddy /usr/bin/caddy

COPY Caddyfile /etc/caddy/Caddyfile
COPY index.md /www/index.md

EXPOSE 2016

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["caddy"]
