FROM alpine:3.24 AS builder

RUN apk add --no-cache g++ linux-headers make

ENV VER=0.745

WORKDIR /luau

ADD https://github.com/luau-lang/luau.git#$VER .

RUN sed -ri 's/^((CXX|LD)FLAGS=)/\1-flto -static /' Makefile \
 && make config=release luau \
 && strip build/release/luau

FROM codegolf/lang-base

COPY --from=0 /luau/build/release/luau /usr/bin/

ENTRYPOINT ["luau"]
