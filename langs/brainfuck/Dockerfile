FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=103be5b

WORKDIR /brainfuck

ADD https://github.com/primo-ppcg/bfci.git?commit=$VER .

RUN cc --std=c99 -O3 -o /usr/bin/brainfuck -s -static bfci.c src/*.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/brainfuck /usr/bin/

ENTRYPOINT ["brainfuck"]
