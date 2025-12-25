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

WORKDIR /_build/stage1/lib/x86_64-linux-ghc-$VER-inplace

RUN find . \( -name '*.a' -or -name '*.hi' -or -name '*.p*hi' \) -delete

RUN strip *.so

COPY haskell.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/haskell -s haskell.c

FROM codegolf/lang-base

ARG VER

ENV GHC=lib/x86_64-linux-ghc-$VER-inplace

COPY --from=0 /lib/ld-musl-*.so.1                             /lib/
COPY --from=0 /usr/bin/gcc /usr/bin/haskell                   /usr/bin/
COPY --from=0 /usr/lib/libffi.so.8                            \
              /usr/lib/libgmp.so                              \
              /usr/lib/libgmp.so.10                           \
              /usr/lib/libncursesw.so.6                       /usr/lib/
COPY --from=0 /_build/stage1/bin/ghc                          /usr/local/bin/
COPY --from=0 /_build/stage1/lib/settings                     /usr/local/lib/
COPY --from=0 /_build/stage1/lib/package.conf.d/package.cache /usr/local/lib/package.conf.d/
COPY --from=0 /_build/stage1/$GHC/libHS*-*_thr-*.so           \
              /_build/stage1/$GHC/libHS*-inplace-*.so         /usr/local/$GHC/
COPY --from=0 /_build/stage1/$GHC/array-*-inplace             /usr/local/$GHC/array-0.5.8.0-inplace
COPY --from=0 /_build/stage1/$GHC/base-*-inplace              /usr/local/$GHC/base-4.21.0.0-inplace
COPY --from=0 /_build/stage1/$GHC/bytestring-*-inplace        /usr/local/$GHC/bytestring-0.12.2.0-inplace
COPY --from=0 /_build/stage1/$GHC/containers-*-inplace        /usr/local/$GHC/containers-0.7-inplace
COPY --from=0 /_build/stage1/$GHC/deepseq-*-inplace           /usr/local/$GHC/deepseq-1.5.1.0-inplace
COPY --from=0 /_build/stage1/$GHC/ghc-bignum-*-inplace        /usr/local/$GHC/ghc-bignum-1.3-inplace
COPY --from=0 /_build/stage1/$GHC/ghc-boot-th-*-inplace       /usr/local/$GHC/ghc-boot-th-$VER-inplace
COPY --from=0 /_build/stage1/$GHC/ghc-internal-*-inplace      /usr/local/$GHC/ghc-internal-9.1202.0-inplace
COPY --from=0 /_build/stage1/$GHC/ghc-prim-*-inplace          /usr/local/$GHC/ghc-prim-0.13.0-inplace
COPY --from=0 /_build/stage1/$GHC/pretty-*-inplace            /usr/local/$GHC/pretty-1.1.3.6-inplace
COPY --from=0 /_build/stage1/$GHC/rts-1.0.2                   /usr/local/$GHC/rts-1.0.2
COPY --from=0 /_build/stage1/$GHC/template-haskell-*-inplace  /usr/local/$GHC/template-haskell-2.23.0.0-inplace

ENTRYPOINT ["haskell"]

CMD ["--version"]
