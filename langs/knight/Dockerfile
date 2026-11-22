FROM alpine:3.24 AS builder

RUN apk add --no-cache clang gcc git make

ENV VER=a2c31e4

WORKDIR /knight

ADD https://github.com/synt7x/knightjit.git?commit=$VER .

RUN make CFLAGS='-O3 -march=x86-64-v3 -DJIT_OFF -static' LDFLAGS='-static'

FROM codegolf/lang-base

COPY --from=0 /knight/target/knight /usr/bin/

ENTRYPOINT ["knight"]
