FROM alpine:3.20 AS builder

ENV VER=1.7.1

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/jqlang/jq/releases/download/jq-$VER/jq-$VER.tar.gz | tar xz --strip-components=1

RUN ./configure --with-oniguruma=builtin \
 && make LDFLAGS='-all-static' -j`nproc`

FROM codegolf/lang-base-no-proc

COPY --from=0 ./jq /usr/bin/

ENTRYPOINT ["jq"]

CMD ["--version"]
