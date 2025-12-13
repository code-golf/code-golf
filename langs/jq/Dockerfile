FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=1.8.1

RUN curl -#L https://github.com/jqlang/jq/releases/download/jq-$VER/jq-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure         \
    --enable-all-static \
    --prefix=/usr       \
 && make install        \
 && strip /usr/bin/jq

FROM codegolf/lang-base

COPY --from=0 /usr/bin/jq /usr/bin/

ENTRYPOINT ["jq"]

CMD ["--version"]
