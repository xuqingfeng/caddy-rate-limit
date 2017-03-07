FROM golang:1.8

RUN go get -v github.com/caddyserver/caddydev \
    && go get -v github.com/mholt/caddy/caddy

WORKDIR /go/src/github.com/xuqingfeng/caddy-rate-limit

COPY . .

EXPOSE 2016

CMD ["caddydev"]