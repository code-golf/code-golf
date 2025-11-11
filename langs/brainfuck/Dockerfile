FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV VER=103be5b

RUN curl -#L https://github.com/primo-ppcg/bfci/tarball/$VER \
  | tar xz --strip-components 1

RUN cc --std=c99 -O3 -o /usr/bin/brainfuck -s -static bfci.c src/*.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/brainfuck /usr/bin/

ENTRYPOINT ["brainfuck"]
