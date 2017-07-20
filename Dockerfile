FROM golang:1.8

RUN git clone https://github.com/xuqingfeng/caddy.git /go/src/github.com/mholt/caddy

WORKDIR /go/src/github.com/xuqingfeng/caddy-rate-limit

COPY . .

RUN cd /go/src/github.com/mholt/caddy/caddy && \
    ./build.bash && \
    cp caddy /go/src/github.com/xuqingfeng/caddy-rate-limit/caddy

EXPOSE 2016

CMD ["./caddy"]