## caddy-rate-limit ![Travis CI](https://img.shields.io/travis/xuqingfeng/caddy-rate-limit/master.svg?style=flat-square)
>a `rate limit` plugin for [caddy](https://caddyserver.com/)

### Syntax

For single resource:

```
ratelimit path rate burst
```

- path is the file or directory to apply `rate limit`

- rate is the limited request in second (r/s) (eg. 1)

- burst is the burst size requester can exceed (eg. 1)


For multiple resources:

```
ratelimit rate burst {
    resources
}
```

- rate is the rate in second (r/s) (eg. 1)

- burst is the burst size requester can exceed (eg. 1)

- resources is a list of files/directories to apply `rate limit`, one per line

### Examples

`ratelimit /r 2 1`

```
ratelimit 2 2 {
    /r1
    /r2
}
```

#### Inspired By

[http://nginx.org/en/docs/http/ngx_http_limit_req_module.html](http://nginx.org/en/docs/http/ngx_http_limit_req_module.html)

[https://github.com/didip/tollbooth](https://github.com/didip/tollbooth)


###### todo

- [ ] fix burst