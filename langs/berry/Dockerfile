FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl python3 readline-dev

ENV VER=97e0fca

RUN curl -#L https://github.com/berry-lang/berry/archive/$VER.tar.gz \
  | tar xz --strip-components 1

COPY berry.c /default/

RUN make && make install

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/local/bin/berry     /usr/bin/

ENTRYPOINT ["berry"]

CMD ["--version"]
