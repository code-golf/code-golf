FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base cabal curl ncurses-dev

ENV VER=1.01

RUN curl https://golfscript.com/nibbles/nibbles-$VER.tgz \
  | tar xz --directory /usr/local/bin --strip-components 1

RUN curl -#L https://downloads.haskell.org/~ghc/8.10.6/ghc-8.10.6-x86_64-alpine3.10-linux-integer-simple.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure  \
 && make install \
 && cabal update \
 && cabal install --lib dlist memoize murmur-hash split

COPY nibbles.patch \
     nibbleswrapper.c /

RUN patch -p 0 < nibbles.patch

RUN ghc -O -package bytestring \
           -package containers \
           -package filepath   \
           -package ghc        \
           -package process    \
           -package template-haskell /usr/local/bin/*.hs

RUN gcc -Wall -Werror -Wextra -o /usr/bin/nibbleswrapper -s -static nibbleswrapper.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/bin/cc               \
              /usr/bin/nibbleswrapper   /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6 /usr/lib/
COPY --from=0 /usr/local/bin            /usr/local/bin

ENTRYPOINT ["nibbleswrapper"]

CMD ["--version"]
