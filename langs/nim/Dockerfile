FROM alpine:edge

RUN apk add -X http://dl-cdn.alpinelinux.org/alpine/edge/testing --no-cache \
    musl-dev nim tcc

ENTRYPOINT ["/usr/bin/nim", "--cc:tcc", "--verbosity:0", "-r", "c", "-"]
