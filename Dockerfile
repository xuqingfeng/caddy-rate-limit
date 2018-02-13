FROM golang:1.9

RUN git clone https://github.com/xuqingfeng/caddy.git /go/src/github.com/mholt/caddy

RUN go get github.com/caddyserver/builds

WORKDIR /go/src/github.com/xuqingfeng/caddy-rate-limit

COPY . .

RUN cd /go/src/github.com/mholt/caddy/caddy && \
    go run build.go && \
    cp caddy /go/src/github.com/xuqingfeng/caddy-rate-limit/caddy

EXPOSE 2016

CMD ["./caddy"]