FROM golang:1.7

WORKDIR /go/src/github.com/xuqingfeng/caddy-rate-limit

RUN go get -v github.com/caddyserver/caddydev
RUN go get -v github.com/mholt/caddy/caddy
COPY . .

EXPOSE 2016

CMD ["caddydev"]