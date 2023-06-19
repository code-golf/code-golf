FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/janet-lang/janet/archive/refs/tags/v1.29.1.tar.gz | tar xz

RUN make -C janet-* install

FROM codegolf/lang-base

COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=builder /usr/local/bin/janet     /usr/bin/

ENTRYPOINT ["janet"]

CMD ["-v"]
