console.info(`
ratelimit 1 2 second {
    whitelist 127.0.0.1/32
    whitelist 1.2.3.4/32
    /static
}
`);