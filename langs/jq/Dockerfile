FROM alpine:3.20 AS builder

ARG RELEASE=1.7.1

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/jqlang/jq/releases/download/jq-$RELEASE/jq-$RELEASE.tar.gz | tar xz

WORKDIR /jq-$RELEASE

RUN ./configure --with-oniguruma=builtin \
 && make -j8 && make check && make install

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/local/bin/jq        /usr/bin/

ENTRYPOINT ["jq"]

CMD ["--version"]
