FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

ENV LDFLAGS='-static' LIBJANET_LDFLAGS='' VER=1.41.2

WORKDIR /janet

ADD https://github.com/janet-lang/janet.git#v$VER .

RUN make build/janet install \
 && strip /usr/local/bin/janet

FROM codegolf/lang-base

COPY --from=0 /usr/local/bin/janet /usr/bin/

ENTRYPOINT ["janet"]

CMD ["--version"]
