FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV CC='gcc -static' VER=4dd7915334

WORKDIR /usr/bin

RUN curl -#L https://codeberg.org/ngn/k/archive/$VER.tar.gz \
  | tar xz --strip-components 1

RUN sed -i 's/native/x86-64-v3/' makefile | make

FROM codegolf/lang-base

COPY --from=0 /usr/bin/k /usr/bin/

ENTRYPOINT ["k"]
