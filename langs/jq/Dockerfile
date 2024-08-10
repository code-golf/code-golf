FROM alpine:3.20 AS builder

ENV VER=1.7.1

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/jqlang/jq/releases/download/jq-$VER/jq-$VER.tar.gz | tar xz

WORKDIR /jq-$VER

RUN ./configure --with-oniguruma=builtin \
 && make -j`nproc` && make check && make install

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/local/bin/jq        /usr/bin/

ENTRYPOINT ["jq"]

CMD ["--version"]
