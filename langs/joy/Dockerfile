FROM alpine:3.22 AS builder

RUN apk add --no-cache bison cmake curl flex g++ gc-dev gc-static make

ENV LDFLAGS='-static' VER=a68b339

RUN curl -#L https://github.com/Wodan58/Moy/tarball/$VER \
  | tar xz --strip-components 1

RUN sed -i 's/(argv\[i]\[0]/(argv[0][0]/' main.c

WORKDIR /usr/bin

RUN cmake ../..     \
 && cmake --build . \
 && strip joy

FROM codegolf/lang-base

COPY --from=0 /usr/bin/joy /usr/bin/
