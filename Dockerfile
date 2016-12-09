FROM golang:1.7

WORKDIR /go/src/github.com/xuqingfeng/caddy-rate-limit

RUN go get -d github.com/xuqingfeng/caddy-rate-limit \
    && go get github.com/caddyserver/caddydev \

CMD ["caddydev"]