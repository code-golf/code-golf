FROM alpine:3.23 AS builder

RUN apk add --no-cache bison cmake curl flex g++ gc-dev gc-static make patch

ENV LDFLAGS='-static' VER=a68b339

RUN curl -#L https://github.com/Wodan58/Moy/tarball/$VER \
  | tar xz --strip-components 1

COPY joy.patch /

RUN patch -p0 < joy.patch

WORKDIR /usr/bin

RUN cmake ../..     \
 && cmake --build . \
 && strip joy

FROM codegolf/lang-base

COPY --from=0 /usr/bin/joy /usr/bin/
