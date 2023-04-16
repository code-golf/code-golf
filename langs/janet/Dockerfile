FROM alpine:3.17 AS builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/janet-lang/janet/archive/refs/tags/v1.27.0.tar.gz | tar xz

RUN mv janet-* janet \
 && cd janet         \
 && make             \
 && strip build/janet

FROM codegolf/lang-base

COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=builder /janet/build/janet /usr/bin/janet

ENTRYPOINT ["janet"]

CMD ["-v"]
