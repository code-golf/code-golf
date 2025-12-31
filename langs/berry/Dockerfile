FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc make musl-dev python3 readline-dev

ENV VER=7af8289

RUN curl -#L https://github.com/berry-lang/berry/archive/$VER.tar.gz \
  | tar xz --strip-components 1

COPY berry.c default/

RUN make && make install

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1  /lib/
COPY --from=0 /usr/local/bin/berry /usr/bin/
