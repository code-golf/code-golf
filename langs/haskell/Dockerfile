ARG VER=9.12.2

FROM alpine:3.23 AS builder

ARG VER

RUN apk add --no-cache bash build-base cabal curl ncurses-dev python3

RUN curl -#L https://downloads.haskell.org/~ghc/$VER/ghc-$VER-src.tar.xz \
  | tar xJ --strip-components 1

RUN cabal update      \
 && ./configure       \
 && hadrian/build -j  \
    --docs=none       \
    --flavour=release \
    stage2:exe:ghc-bin

RUN find _* \( -name '*.a' -or -name '*.hi' -or -name '*.p*hi' \) -delete

COPY haskell.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/haskell -s haskell.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-*.so.1                                      /lib/
COPY --from=0 /usr/bin/gcc /usr/bin/haskell                            /usr/bin/
COPY --from=0 /usr/lib/libffi.so.8                                     \
              /usr/lib/libgmp.so                                       \
              /usr/lib/libgmp.so.10                                    \
              /usr/lib/libncursesw.so.6                                /usr/lib/
COPY --from=0 /_build/stage1/bin/ghc                                   /usr/local/bin/
COPY --from=0 /_build/stage1/lib/settings                              /usr/local/lib/
COPY --from=0 /_build/stage1/lib/package.conf.d/package.cache          /usr/local/lib/package.conf.d/
COPY --from=0 /_build/stage1/lib/*/libHS*-*_thr-*.so                   \
              /_build/stage1/lib/*/libHS*-inplace-*.so                 /usr/local/lib/x86_64-linux-ghc-$VER-inplace/
COPY --from=0 /_build/stage1/lib/*/base-*-inplace                      /usr/local/lib/x86_64-linux-ghc-$VER-inplace/base-4.21.0.0-inplace
COPY --from=0 /_build/stage1/lib/*/ghc-bignum-*-inplace/GHC/Num        /usr/local/lib/x86_64-linux-ghc-$VER-inplace/ghc-bignum-1.3-inplace/GHC/Num
COPY --from=0 /_build/stage1/lib/*/ghc-internal-*-inplace/GHC/Internal /usr/local/lib/x86_64-linux-ghc-$VER-inplace/ghc-internal-9.1202.0-inplace/GHC/Internal
COPY --from=0 /_build/stage1/lib/*/ghc-prim-*-inplace/GHC              /usr/local/lib/x86_64-linux-ghc-$VER-inplace/ghc-prim-0.13.0-inplace/GHC
COPY --from=0 /_build/stage1/lib/*/rts-*/include/DerivedConstants.h    /usr/local/lib/x86_64-linux-ghc-$VER-inplace/rts-1.0.2/include/

ENTRYPOINT ["haskell"]

CMD ["--version"]
