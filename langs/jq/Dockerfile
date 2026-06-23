FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

ENV VER=1.8.2

ADD --unpack https://github.com/jqlang/jq/releases/download/jq-$VER/jq-$VER.tar.gz /

WORKDIR /jq-$VER

RUN ./configure         \
    --enable-all-static \
    --prefix=/usr       \
 && make install        \
 && strip /usr/bin/jq

FROM codegolf/lang-base

COPY --from=0 /usr/bin/jq /usr/bin/

ENTRYPOINT ["jq"]

CMD ["--version"]
