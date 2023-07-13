FROM alpine:3.18 AS builder

RUN apk add -X http://dl-cdn.alpinelinux.org/alpine/edge/testing --no-cache \
    gforth=0.7.3-r3

FROM codegolf/lang-base

COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=builder /usr/bin/gforth          /usr/bin/forth
COPY --from=builder /usr/include/gforth      /usr/include/gforth
COPY --from=builder /usr/lib/gforth          /usr/lib/gforth
COPY --from=builder /usr/lib/libltdl.so.7    /usr/lib/
COPY --from=builder /usr/share/gforth        /usr/share/gforth

ENTRYPOINT ["forth"]

CMD ["--version"]
