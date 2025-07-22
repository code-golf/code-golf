FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV LDFLAGS='-static' VER=1.38.0

RUN curl -#L https://github.com/janet-lang/janet/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN make build/janet install \
 && strip /usr/local/bin/janet

FROM codegolf/lang-base

COPY --from=0 /usr/local/bin/janet /usr/bin/

ENTRYPOINT ["janet"]

CMD ["--version"]
