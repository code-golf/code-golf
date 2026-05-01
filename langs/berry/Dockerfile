FROM alpine:3.23 AS builder

RUN apk add --no-cache gcc make musl-dev python3 readline-static

ENV VER=26e9cb5

WORKDIR /berry

ADD https://github.com/berry-lang/berry.git?commit=$VER .

COPY berry.c default/

RUN make CFLAGS='-flto -O2 -static' LFLAGS='-flto -O2 -static'

FROM codegolf/lang-base

COPY --from=0 /berry/berry /usr/bin/
