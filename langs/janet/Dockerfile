FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/janet-lang/janet/archive/refs/tags/v1.30.0.tar.gz | tar xz

RUN mv janet-* janet                          \
 && make -C janet build/janet LDFLAGS=-static \
 && strip janet/build/janet

FROM codegolf/lang-base

COPY --from=builder /janet/build/janet /usr/bin/

ENTRYPOINT ["janet"]

CMD ["-v"]
