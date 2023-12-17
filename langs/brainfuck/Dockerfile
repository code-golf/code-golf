FROM alpine:3.18 AS builder

ENV VERSION=874e785

RUN apk add --no-cache curl gcc musl-dev

RUN curl -L https://github.com/primo-ppcg/bfci/tarball/$VERSION | tar xz \
 && mv primo-ppcg-bfci-* bfci

RUN cd bfci && cc --std=c99 -static -O3 -s -o bfci bfci.c src/*.c

FROM codegolf/lang-base

COPY --from=0 /bfci/bfci /usr/bin/brainfuck

ENTRYPOINT ["brainfuck"]
