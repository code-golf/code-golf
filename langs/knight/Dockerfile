FROM alpine:3.23 AS builder

RUN apk add --no-cache clang git make

ENV VER=1981682

WORKDIR /knight

ADD https://github.com/synt7x/knightjit.git?commit=$VER /knight

RUN make CFLAGS='-O3 -march=x86-64-v3 -DJIT_OFF -static' LDFLAGS='-static'

FROM codegolf/lang-base

COPY --from=0 /knight/target/knight /usr/bin/

ENTRYPOINT ["knight"]
