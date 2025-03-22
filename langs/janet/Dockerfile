FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl

ENV VER=1.38.0

RUN curl -#L https://github.com/janet-lang/janet/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN LDFLAGS='-static' make build/janet \
 && strip /build/janet

FROM codegolf/lang-base

COPY --from=0 /build/janet /usr/bin/

ENTRYPOINT ["janet"]

CMD ["--version"]
