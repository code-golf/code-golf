FROM alpine:3.23 AS builder

ENV GHCVER=9.8.2 VER=1.1.4

RUN apk add --no-cache ghc=$GHCVER-r2 gmp-static musl-dev

ADD https://golfscript.com/iogii/$VER/iogii /usr/local/bin/iogii

ADD https://golfscript.com/iogii/$VER/engine.hs /usr/bin/ghcii.hs

RUN sed -i '1c#!/usr/bin/ruby --disable=gems' /usr/local/bin/iogii

RUN ghc -O2 -optl -static /usr/bin/ghcii.hs && strip /usr/bin/ghcii

COPY iogii.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iogii -s iogii.c

FROM codegolf/lang-base

COPY --from=codegolf/lang-ruby --parents /lib/ld-musl-*.so.1  \
                                         /usr/bin/ruby        /
COPY --from=0 --chmod=755      --parents /usr/bin/ghcii       \
                                         /usr/bin/iogii       \
                                         /usr/local/bin/iogii /

ENTRYPOINT ["iogii"]

CMD ["--version"]
